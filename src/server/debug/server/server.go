package server

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/coverage"
	runtimedebug "runtime/debug"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wcharczuk/go-chart"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	describe "k8s.io/kubectl/pkg/describe"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"

	"github.com/pachyderm/pachyderm/v2/src/client"
	"github.com/pachyderm/pachyderm/v2/src/debug"
	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
	"github.com/pachyderm/pachyderm/v2/src/internal/errutil"
	"github.com/pachyderm/pachyderm/v2/src/internal/grpcutil"
	"github.com/pachyderm/pachyderm/v2/src/internal/log"
	loki "github.com/pachyderm/pachyderm/v2/src/internal/lokiutil/client"
	"github.com/pachyderm/pachyderm/v2/src/internal/pachsql"
	"github.com/pachyderm/pachyderm/v2/src/internal/ppsutil"
	"github.com/pachyderm/pachyderm/v2/src/internal/serviceenv"
	"github.com/pachyderm/pachyderm/v2/src/internal/uuid"
	"github.com/pachyderm/pachyderm/v2/src/pfs"
	"github.com/pachyderm/pachyderm/v2/src/pps"
)

const (
	defaultDuration = time.Minute
	pachdPrefix     = "pachd"
	pipelinePrefix  = "pipelines"
	podPrefix       = "pods"
	databasePrefix  = "database"
)

type TaskType int

const (
	TaskType_Helm TaskType = iota
	TaskType_SourceRepos
	TaskType_PachdDescribe
	TaskType_PachdLogs
	TaskType_PachdLokiLogs
	TaskType_Database
	TaskType_Version
	TaskType_Profile
	TaskType_AppLogs
	TaskType_Pipelines
	TaskType_Pipeline
	TaskType_LokiLogs
	TaskType_Logs
)

// maps the task type to its sub-directory in the debug dump
func (t TaskType) String() string {
	return []string{
		"helm",
		"source-repos",
		"pachd-describe",
		"pachd-logs",
		"pachd-loki-logs",
		"database",
		"version",
		"profile",
		"app-logs",
		"pipelines",
		"pipeline",
	}[t]
}

type debugServer struct {
	env                 serviceenv.ServiceEnv
	name                string
	sidecarClient       *client.APIClient
	marshaller          *jsonpb.Marshaler
	database            *pachsql.DB
	logLevel, grpcLevel log.LevelChanger
}

// NewDebugServer creates a new server that serves the debug api over GRPC
func NewDebugServer(env serviceenv.ServiceEnv, name string, sidecarClient *client.APIClient, db *pachsql.DB) debug.DebugServer {
	return &debugServer{
		env:           env,
		name:          name,
		sidecarClient: sidecarClient,
		marshaller:    &jsonpb.Marshaler{Indent: "  "},
		database:      db,
		logLevel:      log.LogLevel,
		grpcLevel:     log.GRPCLevel,
	}
}

type collectPipelineFunc func(context.Context, *tar.Writer, *client.APIClient, *pps.PipelineInfo, ...string) error
type collectWorkerFunc func(context.Context, *tar.Writer, *v1.Pod, ...string) error
type redirectFunc func(context.Context, debug.DebugClient, *debug.Filter) (io.Reader, error)
type collectFunc func(context.Context, *tar.Writer, *client.APIClient, ...string) error

type task struct {
	taskType TaskType
	timeout  time.Duration
	app      string
}

// the returned taskPath gets deleted by the caller after its contents are streamed, and is the location of any associated error file
func (s *debugServer) runTask(ctx context.Context, taskType TaskType, dir string, reportProgress func(progress, total int)) (taskPath string, retErr error) {
	ctx, end := log.SpanContext(ctx, "collect_"+taskType.String(), zap.String("path", dir))
	defer func() {
		retErr = errors.Wrap(retErr, " collect_"+taskType.String())
		end(log.Errorp(&retErr))
		if retErr != nil {
			if err := writeErrorFileV2(taskPath, retErr); err != nil {
				retErr = errors.Wrapf(err, "write errors file to dir %q", dir)
			}
		}
	}()
	switch taskType {
	case TaskType_Helm:
		return s.collectHelm(ctx, dir, reportProgress)
	case TaskType_SourceRepos:
		return s.collectInputRepos(ctx, s.env.GetPachClient(ctx), filepath.Join(dir, taskType.String()), int64(0), reportProgress)
	case TaskType_PachdDescribe:
		return s.collectDescribe(ctx, filepath.Join(dir, "pachd"), s.name, reportProgress)
	case TaskType_PachdLogs:
		return s.collectLogs(ctx, filepath.Join(dir, "pachd"), s.name, "pachd", reportProgress)
	case TaskType_PachdLokiLogs:
		return s.collectLogsLoki(ctx, filepath.Join(dir, "pachd"), s.name, "pachd", reportProgress)
	case TaskType_Database:
		return s.collectDatabaseDump(ctx, filepath.Join(dir, databasePrefix), reportProgress)
	case TaskType_Version:
		return s.collectPachdVersion(ctx, s.env.GetPachClient(ctx), dir, reportProgress)
	case TaskType_Profile:
		return collectDump(ctx, dir, reportProgress)
	case TaskType_Pipelines:
		pis, err := s.env.GetPachClient(ctx).ListPipeline(true)
		if err != nil {
			return filepath.Join(dir, pipelinePrefix), errors.Wrap(err, "list pipelines")
		}
		for _, pi := range pis {
			dir := filepath.Join(dir, pipelinePrefix, pi.Pipeline.Project.Name, pi.Pipeline.Name)
			defer func() {
				if retErr != nil {
					retErr = writeErrorFileV2(dir, retErr)
				}
			}()
			// is this necessary????
			if err := validatePipelineInfo(pi); err != nil {
				return filepath.Join(dir, pipelinePrefix), errors.Wrap(err, "collectPipelineDumpFunc: invalid pipeline info")
			}
			if err := s.collectPipelineV2(ctx, dir, pi.Pipeline, 0); err != nil {
				return filepath.Join(dir, pipelinePrefix), err
			}
			// return s.forEachWorker(ctx, pi, func(pod *v1.Pod) error {
			// 	return s.handleWorkerRedirect(ctx, tw, pod, collectWorker, redirect, prefix)
			// })
		}
		return filepath.Join(dir, pipelinePrefix), nil
	case TaskType_Pipeline: // TODO(acohen4): fill this in
		p := &pps.Pipeline{}
		pi, err := s.env.GetPachClient(ctx).InspectProjectPipeline(p.Project.GetName(), p.Name, true)
		if err != nil {
			return "", errors.Wrapf(err, "inspect pipeline %q", p.String())
		}
		if err := validatePipelineInfo(pi); err != nil {
			return "", errors.Wrap(err, "collectPipelineDumpFunc: invalid pipeline info")
		}
		dir := filepath.Join(dir, pipelinePrefix, p.Project.Name, p.Name)
		defer func() {
			if retErr != nil {
				retErr = writeErrorFileV2(dir, retErr)
			}
		}()
		if err := s.collectPipelineV2(ctx, dir, pi.Pipeline, 0); err != nil {
			return filepath.Join(dir, pipelinePrefix), err
		}
		// return s.forEachWorker(ctx, pi, func(pod *v1.Pod) error {
		// 	return s.handleWorkerRedirect(ctx, tw, pod, collectWorker, redirect, prefix)
		// })
		return filepath.Join(dir, pipelinePrefix), nil
	case TaskType_AppLogs: // TODO: what the heck is this anyway?
		return "", s.appLogs(ctx, dir)
	default:
		return "", errors.Errorf("no task implemetation for %q", taskType)
	}
}

func (s *debugServer) GetDumpV2Template(ctx context.Context, request *debug.GetDumpV2TemplateRequest) (*debug.GetDumpV2TemplateResponse, error) {
	var ps []string
	var pipTasks []*debug.PipelineDump
	for _, p := range ps {
		pipTasks = append(pipTasks, &debug.PipelineDump{Name: p})
	}
	return &debug.GetDumpV2TemplateResponse{
		Request: &debug.DumpV2Request{
			SystemDump: &debug.SystemDump{
				Version: true,
				Helm:    true,
				Profile: true,
			},
			PachdDump: &debug.PachdDump{
				Describe: true,
				Logs:     true, // do we want to either do Logs or LokiLogs based on env?
				LokiLogs: true,
			},
			AppDump: &debug.AppDump{
				InputRepos: true,
			},
			PipelinesDump: &debug.PipelinesDump{
				Tasks: pipTasks,
			},
		},
	}, nil
}

func (s *debugServer) DumpV2(request *debug.DumpV2Request, server debug.Debug_DumpV2Server) error {
	return s.dump(s.env.GetPachClient(server.Context()), server, dumpTasks(request))
}

func dumpTasks(request *debug.DumpV2Request) []task {
	var tasks []task
	if request.AppDump != nil {
		if request.AppDump.InputRepos {
			tasks = append(tasks, task{taskType: TaskType_SourceRepos})
		}
	}
	if request.PachdDump != nil {
		if request.PachdDump.Describe {
			tasks = append(tasks, task{taskType: TaskType_PachdDescribe})
		}
		if request.PachdDump.Logs {
			tasks = append(tasks, task{
				taskType: TaskType_PachdLogs,
				timeout:  10 * time.Minute,
			})
		}
		if request.PachdDump.LokiLogs {
			tasks = append(tasks, task{
				taskType: TaskType_PachdLokiLogs,
				timeout:  10 * time.Minute,
			})
		}
	}
	if request.PipelinesDump != nil {
		for _, t := range request.PipelinesDump.Tasks {
			if len(t.Tasks) == 0 { // expand to all pipelines

			} else {

			}
		}
	}
	if request.SystemDump != nil {
		if request.SystemDump.Version {
			tasks = append(tasks, task{taskType: TaskType_Version})
		}
		if request.SystemDump.Helm {
			tasks = append(tasks, task{taskType: TaskType_Helm})
		}
		if request.SystemDump.Profile {
			tasks = append(tasks, task{taskType: TaskType_Profile})
		}
	}
	return tasks
}

type dumpContentServer struct {
	server debug.Debug_DumpV2Server
}

func (dcs *dumpContentServer) Send(bytesValue *types.BytesValue) error {
	return dcs.server.Send(&debug.DumpChunk{Chunk: &debug.DumpChunk_Content{Content: &debug.DumpContent{Content: bytesValue.Value}}})
}

func sendProgressFunc(ctx context.Context, server debug.Debug_DumpV2Server, t task) func(progress, total int) {
	return func(progress, total int) {
		if err := server.Send(
			&debug.DumpChunk{
				Chunk: &debug.DumpChunk_Progress{
					Progress: &debug.DumpProgress{
						Task:     t.taskType.String(),
						Progress: int64(progress),
						Total:    int64(total),
					},
				},
			},
		); err != nil {
			log.Error(ctx, fmt.Sprintf("error sending progress for task %q", t.taskType), zap.Error(err))
		}
	}

}

func (s *debugServer) dump(
	c *client.APIClient,
	server debug.Debug_DumpV2Server,
	tasks []task,
) (retErr error) {
	ctx := c.Ctx() // this context has authorization credentials we need
	dumpRoot := filepath.Join(os.TempDir(), uuid.NewWithoutDashes())
	if err := os.Mkdir(dumpRoot, os.ModeDir); err != nil {
		return errors.Wrap(err, "create dump root directory")
	}
	defer func() {
		retErr = multierr.Append(retErr, os.RemoveAll(dumpRoot))
	}()
	eg, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex // to synchronize writes to the tar stream
	dumpContent := &dumpContentServer{server: server}
	return grpcutil.WithStreamingBytesWriter(dumpContent, func(w io.Writer) error {
		return withDebugWriter(w, func(tw *tar.Writer) error {
			for _, t := range tasks {
				func(ctx context.Context, t task) {
					eg.Go(func() (retErr error) {
						taskCtx := ctx
						defer func() {
							if retErr != nil {
								log.Error(taskCtx, fmt.Sprintf("failed dump task %q", t.taskType.String()), zap.Error(retErr))
								retErr = nil
							}
						}()
						if t.timeout != 0 {
							var cf context.CancelFunc
							ctx, cf = context.WithTimeout(ctx, t.timeout)
							defer cf()
						}
						outPath, err := s.runTask(ctx, t.taskType, dumpRoot, sendProgressFunc(ctx, server, t))
						if err != nil {
							return errors.Wrapf(err, "run task %q", t.taskType.String())
						}
						if outPath == "" {
							return
						}
						defer func() {
							if err := os.RemoveAll(outPath); err != nil {
								retErr = multierr.Append(retErr, err)
							}
						}()
						mu.Lock()
						defer mu.Unlock()
						return filepath.Walk(outPath, func(path string, fi fs.FileInfo, err error) error {
							if err != nil {
								return errors.Wrapf(err, "walk path %q for dump task %q", path, t.taskType.String())
							}
							if fi.IsDir() {
								return nil
							}
							dest := strings.TrimPrefix(path, dumpRoot)
							return errors.Wrapf(writeTarFileV2(tw, dest, path, fi), "write tar file %q", dest)
						})
					})
				}(ctx, t)
			}
			return errors.Wrap(eg.Wait(), "run dump tasks")
		})
	})
}

func (s *debugServer) appLogs(ctx context.Context, dir string) (retErr error) {
	ctx, end := log.SpanContext(ctx, "collectAppLogs")
	defer end(log.Errorp(&retErr))
	ctx, c := context.WithTimeout(ctx, 10*time.Minute)
	defer c()
	pods, err := s.env.GetKubeClient().CoreV1().Pods(s.env.Config().Namespace).List(ctx, metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ListOptions",
			APIVersion: "v1",
		},
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: map[string]string{
				"suite": "pachyderm",
			},
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "component",
					Operator: metav1.LabelSelectorOpNotIn,
					// Worker and pachd logs are collected by separate
					// functions.
					Values: []string{"worker", "pachd"},
				},
			},
		}),
	})
	if err != nil {
		return errors.EnsureStack(err)
	}
	var errs error
	for _, pod := range pods.Items {
		dir := filepath.Join(dir, pod.Labels["app"], pod.Name)
		if _, err := s.collectDescribe(ctx, dir, pod.Name, func(_, _ int) {}); err != nil {
			multierr.AppendInto(&errs, errors.Wrapf(err, "describe(%s)", pod.Name))
		}
		for _, container := range pod.Spec.Containers {
			dir := filepath.Join(dir, container.Name)
			if _, err := s.collectLogs(ctx, dir, pod.Name, container.Name, func(progress, total int) {}); err != nil {
				multierr.AppendInto(&errs, errors.Wrapf(err, "collectLogs(%s.%s)", pod.Name, container.Name))
			}
			if _, err := s.collectLogsLoki(ctx, dir, pod.Name, container.Name, func(progress, total int) {}); err != nil {
				multierr.AppendInto(&errs, errors.Wrapf(err, "collectLogsLoki(%s.%s)", pod.Name, container.Name))
			}
		}
	}
	return errs
}

func (s *debugServer) forEachWorker(ctx context.Context, pipelineInfo *pps.PipelineInfo, cb func(*v1.Pod) error) error {
	pods, err := s.getWorkerPods(ctx, pipelineInfo)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		return errors.Errorf("no worker pods found for pipeline %s", pipelineInfo.GetPipeline())
	}
	var errs error
	for _, pod := range pods {
		if err := cb(&pod); err != nil {
			multierr.AppendInto(&errs, errors.Wrapf(err, "forEachWorker(%s)", pod.Name))
		}
	}
	return nil
}

func (s *debugServer) getWorkerPods(ctx context.Context, pipelineInfo *pps.PipelineInfo) ([]v1.Pod, error) {
	podList, err := s.env.GetKubeClient().CoreV1().Pods(s.env.Config().Namespace).List(
		ctx,
		metav1.ListOptions{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ListOptions",
				APIVersion: "v1",
			},
			LabelSelector: metav1.FormatLabelSelector(
				metav1.SetAsLabelSelector(
					map[string]string{
						"app":             "pipeline",
						"pipelineProject": pipelineInfo.Pipeline.Project.Name,
						"pipelineName":    pipelineInfo.Pipeline.Name,
						"pipelineVersion": fmt.Sprint(pipelineInfo.Version),
					},
				),
			),
		},
	)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return s.getLegacyWorkerPods(ctx, pipelineInfo)
		}
		return nil, errors.EnsureStack(err)
	}
	if len(podList.Items) == 0 {
		return s.getLegacyWorkerPods(ctx, pipelineInfo)
	}
	return podList.Items, nil
}

// DELETE THIS
func (s *debugServer) getLegacyWorkerPods(ctx context.Context, pipelineInfo *pps.PipelineInfo) ([]v1.Pod, error) {
	podList, err := s.env.GetKubeClient().CoreV1().Pods(s.env.Config().Namespace).List(
		ctx,
		metav1.ListOptions{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ListOptions",
				APIVersion: "v1",
			},
			LabelSelector: metav1.FormatLabelSelector(
				metav1.SetAsLabelSelector(
					map[string]string{
						"app": ppsutil.PipelineRcName(pipelineInfo),
					},
				),
			),
		},
	)
	if err != nil {
		return nil, errors.EnsureStack(err)
	}
	return podList.Items, nil
}

func (s *debugServer) collectWorker(ctx context.Context, dir string, pod *v1.Pod) (retErr error) {
	ctx, end := log.SpanContext(ctx, "collectWorker", zap.String("pod", pod.Name))
	defer end(log.Errorp(&retErr))
	workerDir := filepath.Join(dir, podPrefix, pod.Name)
	defer func() {
		if retErr != nil {
			retErr = writeErrorFileV2(workerDir, retErr)
		}
	}()
	return nil
}

func (s *debugServer) Profile(request *debug.ProfileRequest, server debug.Debug_ProfileServer) error {
	//pachClient := s.env.GetPachClient(server.Context())
	return errors.New("TODO")
	// return s.handleRedirect(
	// 	pachClient,
	// 	server,
	// 	request.Filter,
	// 	collectProfileFunc(request.Profile),  /* collectPachd */
	// 	nil,                                  /* collectPipeline */
	// 	nil,                                  /* collectWorker */
	// 	redirectProfileFunc(request.Profile), /* redirect */
	// 	collectProfileFunc(request.Profile),  /* collect */
	// 	false,
	// 	false,
	// )
}

func collectProfileFunc(profile *debug.Profile) collectFunc {
	return func(ctx context.Context, tw *tar.Writer, _ *client.APIClient, prefix ...string) error {
		switch profile.GetName() {
		case "cover":
			var errs error
			// "go tool covdata" relies on these files being named the same as what is
			// written out automatically.  See go/src/runtime/coverage/emit.go.
			// (covcounters.<16 byte hex-encoded hash that nothing checks>.<pid>.<unix
			// nanos>, covmeta.<same hash>.)
			hash := [16]byte{}
			if _, err := rand.Read(hash[:]); err != nil {
				return errors.Wrap(err, "generate random coverage id")
			}
			if err := collectDebugFile(tw, fmt.Sprintf("cover/covcounters.%x.%d.%d", hash, os.Getpid(), time.Now().UnixNano()), "", coverage.WriteCounters, prefix...); err != nil {
				multierr.AppendInto(&errs, errors.Wrap(err, "counters"))
			}
			if err := collectDebugFile(tw, fmt.Sprintf("cover/covmeta.%x", hash), "", coverage.WriteMeta, prefix...); err != nil {
				multierr.AppendInto(&errs, errors.Wrap(err, "meta"))
			}
			if errs == nil {
				log.Debug(ctx, "clearing coverage counters")
				if err := coverage.ClearCounters(); err != nil {
					log.Debug(ctx, "problem clearing coverage counters", zap.Error(err))
				}
			}
			return errs
		default:
			//return collectProfile(ctx, tw, profile, prefix...)
			return nil
		}

	}
}

func collectProfile(ctx context.Context, dir string, profile *debug.Profile) error {
	return collectDebugFileV2(dir, profile.Name, func(w io.Writer) error {
		return writeProfile(ctx, w, profile)
	})
}

func writeProfile(ctx context.Context, w io.Writer, profile *debug.Profile) (retErr error) {
	defer log.Span(ctx, "writeProfile", zap.String("profile", profile.GetName()))(log.Errorp(&retErr))
	if profile.Name == "cpu" {
		if err := pprof.StartCPUProfile(w); err != nil {
			return errors.EnsureStack(err)
		}
		defer pprof.StopCPUProfile()
		duration := defaultDuration
		if profile.Duration != nil {
			var err error
			duration, err = types.DurationFromProto(profile.Duration)
			if err != nil {
				return errors.EnsureStack(err)
			}
		}
		// Wait for either the defined duration, or until the context is
		// done.
		t := time.NewTimer(duration)
		select {
		case <-ctx.Done():
			t.Stop()
			return errors.EnsureStack(ctx.Err())
		case <-t.C:
			return nil
		}
	}
	p := pprof.Lookup(profile.Name)
	if p == nil {
		return errors.Errorf("unable to find profile %q", profile.Name)
	}
	return errors.EnsureStack(p.WriteTo(w, 0))
}

func (s *debugServer) Binary(request *debug.BinaryRequest, server debug.Debug_BinaryServer) error {
	// pachClient := s.env.GetPachClient(server.Context())
	// return s.handleRedirect(
	// 	pachClient,
	// 	server,
	// 	request.Filter,
	// 	collectBinary,  /* collectPachd */
	// 	nil,            /* collectPipeline */
	// 	nil,            /* collectWorker */
	// 	redirectBinary, /* redirect */
	// 	collectBinary,  /* collect */
	// 	false,
	// 	false,
	// )
	return errors.New("TODO")
}

func collectBinary(ctx context.Context, tw *tar.Writer, _ *client.APIClient, prefix ...string) error {
	return collectDebugFile(tw, "binary", "", func(w io.Writer) (retErr error) {
		defer log.Span(ctx, "collectBinary", zap.String("binary", os.Args[0]))(log.Errorp(&retErr))
		f, err := os.Open(os.Args[0])
		if err != nil {
			return errors.EnsureStack(err)
		}
		defer func() {
			if err := f.Close(); retErr == nil {
				retErr = err
			}
		}()
		_, err = io.Copy(w, f)
		return errors.EnsureStack(err)
	}, prefix...)
}

func redirectBinary(ctx context.Context, c debug.DebugClient, filter *debug.Filter) (io.Reader, error) {
	binaryC, err := c.Binary(ctx, &debug.BinaryRequest{Filter: filter})
	if err != nil {
		return nil, errors.EnsureStack(err)
	}
	return grpcutil.NewStreamingBytesReader(binaryC, nil), nil
}

func (s *debugServer) Dump(request *debug.DumpRequest, server debug.Debug_DumpServer) error {
	return errors.New("Dump gRPC is deprecated in favor of DumpV2")
	// if request.Limit == 0 {
	// 	request.Limit = math.MaxInt64
	// }
	// timeout := 25 * time.Minute // If no client-side deadline set, use 25 minutes.
	// if deadline, ok := server.Context().Deadline(); ok {
	// 	d := time.Until(deadline)
	// 	// Time out our own operations at ~90% of the client-set deadline. (22 minutes of
	// 	// server-side processing for every 25 minutes of debug dumping.)
	// 	timeout = time.Duration(22) * d / time.Duration(25)
	// }
	// ctx, c := context.WithTimeout(server.Context(), timeout)
	// defer c()
	// pachClient := s.env.GetPachClient(ctx)
	// return s.handleRedirect(
	// 	pachClient,
	// 	server,
	// 	request.Filter,
	// 	s.collectPachdDumpFunc(request.Limit), /* collectPachd */
	// 	nil,                                   /* collectPipeline */
	// 	s.collectWorkerDump,                   /* collectWorker */
	// 	redirectDump,                          /* redirect */
	// 	collectDump,                           /* collect */
	// 	true,
	// 	true,
	// )
}

func (s *debugServer) collectPachdDumpFunc(limit int64) collectFunc {
	return func(ctx context.Context, tw *tar.Writer, pachClient *client.APIClient, prefix ...string) (retErr error) {
		return nil
	}
}

func (s *debugServer) collectInputRepos(ctx context.Context, pachClient *client.APIClient, dir string, limit int64, reportProgress func(progress, total int)) (_ string, retErr error) {
	if err := os.Mkdir(dir, os.ModeDir); err != nil {
		return "", errors.Wrap(err, "make source-repos dir")
	}
	reportProgress(0, 1)
	repoInfos, err := pachClient.ListRepo()
	if err != nil {
		return "", errors.Wrap(err, "list repo")
	}
	var errs error
	for i, repoInfo := range repoInfos {
		if err := validateRepoInfo(repoInfo); err != nil {
			multierr.AppendInto(&errs, errors.Wrapf(err, "invalid repo info %d (%s) from ListRepo", i, repoInfo.String()))
			continue
		}
		if _, err := pachClient.InspectProjectPipeline(repoInfo.Repo.Project.Name, repoInfo.Repo.Name, true); err != nil {
			if errutil.IsNotFoundError(err) {
				dir := filepath.Join(dir, repoInfo.Repo.Project.Name, repoInfo.Repo.Name)
				if err := s.collectCommits(ctx, pachClient, dir, repoInfo.Repo, limit); err != nil {
					multierr.AppendInto(&errs, errors.Wrapf(err, "collectCommits(%s:%s)", repoInfo.GetRepo(), dir))
				}
				continue
			}
			multierr.AppendInto(&errs, errors.Wrapf(err, "inspectPipeline(%s)", repoInfo.GetRepo()))
		}
		reportProgress(i, len(repoInfos))
	}
	reportProgress(1, 1)
	return dir, errs
}

func (s *debugServer) collectCommits(rctx context.Context, pachClient *client.APIClient, dir string, repo *pfs.Repo, limit int64) error {
	compacting := chart.ContinuousSeries{
		Name: "compacting",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(255),
		},
	}
	validating := chart.ContinuousSeries{
		Name: "validating",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(1).WithAlpha(255),
		},
	}
	finishing := chart.ContinuousSeries{
		Name: "finishing",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(2).WithAlpha(255),
		},
	}
	if err := collectDebugFileV2(dir, "commits.json", func(w io.Writer) error {
		ctx, cancel := context.WithCancel(rctx)
		defer cancel()
		client, err := pachClient.PfsAPIClient.ListCommit(ctx, &pfs.ListCommitRequest{
			Repo:   repo,
			Number: limit,
			All:    true,
		})
		if err != nil {
			return errors.EnsureStack(err)
		}
		return grpcutil.ForEach[*pfs.CommitInfo](client, func(ci *pfs.CommitInfo) error {
			if ci.Finished != nil && ci.Details != nil && ci.Details.CompactingTime != nil && ci.Details.ValidatingTime != nil {
				compactingDuration, err := types.DurationFromProto(ci.Details.CompactingTime)
				if err != nil {
					return errors.EnsureStack(err)
				}
				compacting.XValues = append(compacting.XValues, float64(len(compacting.XValues)+1))
				compacting.YValues = append(compacting.YValues, float64(compactingDuration))
				validatingDuration, err := types.DurationFromProto(ci.Details.ValidatingTime)
				if err != nil {
					return errors.EnsureStack(err)
				}
				validating.XValues = append(validating.XValues, float64(len(validating.XValues)+1))
				validating.YValues = append(validating.YValues, float64(validatingDuration))
				finishingTime, err := types.TimestampFromProto(ci.Finishing)
				if err != nil {
					return errors.EnsureStack(err)
				}
				finishedTime, err := types.TimestampFromProto(ci.Finished)
				if err != nil {
					return errors.EnsureStack(err)
				}
				finishingDuration := finishedTime.Sub(finishingTime)
				finishing.XValues = append(finishing.XValues, float64(len(finishing.XValues)+1))
				finishing.YValues = append(finishing.YValues, float64(finishingDuration))
			}
			return errors.EnsureStack(s.marshaller.Marshal(w, ci))
		})
	}); err != nil {
		return err
	}
	// Reverse the x values since we collect them from newest to oldest.
	// TODO: It would probably be better to support listing jobs in reverse order.
	reverseContinuousSeries(compacting, validating, finishing)
	return collectGraph(dir, "commits-chart.png", "number of commits", []chart.Series{compacting, validating, finishing})
}

func reverseContinuousSeries(series ...chart.ContinuousSeries) {
	for _, s := range series {
		for i := 0; i < len(s.XValues)/2; i++ {
			s.XValues[i], s.XValues[len(s.XValues)-1-i] = s.XValues[len(s.XValues)-1-i], s.XValues[i]
		}
	}
}

func collectGraph(dir, name, XAxisName string, series []chart.Series, prefix ...string) error {
	return collectDebugFileV2(dir, name, func(w io.Writer) error {
		graph := chart.Chart{
			Title: name,
			TitleStyle: chart.Style{
				Show: true,
			},
			XAxis: chart.XAxis{
				Name: XAxisName,
				NameStyle: chart.Style{
					Show: true,
				},
				Style: chart.Style{
					Show: true,
				},
				ValueFormatter: func(v interface{}) string {
					return fmt.Sprintf("%.0f", v.(float64))
				},
			},
			YAxis: chart.YAxis{
				Name: "time (seconds)",
				NameStyle: chart.Style{
					Show: true,
				},
				Style: chart.Style{
					Show:                true,
					TextHorizontalAlign: chart.TextHorizontalAlignLeft,
				},
				ValueFormatter: func(v interface{}) string {
					seconds := v.(float64) / float64(time.Second)
					return fmt.Sprintf("%.3f", seconds)
				},
			},
			Series: series,
		}
		graph.Elements = []chart.Renderable{
			chart.Legend(&graph),
		}
		return errors.EnsureStack(graph.Render(chart.PNG, w))
	})
}

func collectGoInfo(dir string) error {
	return collectDebugFileV2(dir, "go_info.txt", func(w io.Writer) error {
		fmt.Fprintf(w, "build info: ")
		info, ok := runtimedebug.ReadBuildInfo()
		if ok {
			fmt.Fprintf(w, "%s", info.String())
		} else {
			fmt.Fprint(w, "<no build info>")
		}
		fmt.Fprintf(w, "GOOS: %v\nGOARCH: %v\nGOMAXPROCS: %v\nNumCPU: %v\n", runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0), runtime.NumCPU())
		return nil
	})
}

func (s *debugServer) collectPachdVersion(ctx context.Context, pachClient *client.APIClient, dir string, reportProgress func(progress, total int)) (string, error) {
	defer reportProgress(100, 100)
	err := collectDebugFileV2(dir, "version.txt", func(w io.Writer) (retErr error) {
		version, err := pachClient.Version()
		if err != nil {
			return err
		}
		_, err = io.Copy(w, strings.NewReader(version+"\n"))
		return errors.EnsureStack(err)
	})
	return filepath.Join(dir, "version.txt"), err
}

func (s *debugServer) collectHelm(ctx context.Context, dir string, reportProgress func(progress, total int)) (string, error) {
	secrets, err := s.env.GetKubeClient().CoreV1().Secrets(s.env.Config().Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "owner=helm",
	})
	if err != nil {
		return filepath.Join(dir, "helm"), errors.EnsureStack(err)
	}
	var writeErrs error
	for i, secret := range secrets.Items {
		helmPath := filepath.Join(dir, "helm", secret.Name)
		if err := handleHelmSecretV2(ctx, helmPath, secret); err != nil {
			if err := writeErrorFileV2(helmPath, err); err != nil {
				multierr.AppendInto(&writeErrs, errors.Wrapf(err, "%v: write error.txt", secret.Name))
				continue
			}
		}
		reportProgress(i, len(secrets.Items))
	}
	reportProgress(100, 100)
	return filepath.Join(dir, "helm"), writeErrs
}

func (s *debugServer) collectDescribe(ctx context.Context, dir string, pod string, reportProgress func(progress, total int)) (string, error) {
	err := collectDebugFileV2(dir, "describe.txt", func(output io.Writer) (retErr error) {
		// Gate the total time of "describe".
		ctx, c := context.WithTimeout(ctx, 2*time.Minute)
		defer c()
		defer log.Span(ctx, "collectDescribe")(log.Errorp(&retErr))
		r, w := io.Pipe()
		go func() {
			defer log.Span(ctx, "collectDescribe.backgroundDescribe")()
			// Do the "describe" in the background, because the k8s client doesn't take
			// a context here, and we want the ability to abandon the request.  We leak
			// memory if this runs forever but we return, but that's better than a debug
			// dump that doesn't return.
			pd := describe.PodDescriber{
				Interface: s.env.GetKubeClient(),
			}
			output, err := pd.Describe(s.env.Config().Namespace, pod, describe.DescriberSettings{ShowEvents: true})
			if err != nil {
				w.CloseWithError(errors.EnsureStack(err))
				return
			}
			_, err = w.Write([]byte(output))
			w.CloseWithError(errors.EnsureStack(err))
		}()
		go func() {
			// Close the pipe when the context times out; bounding the time on the
			// io.Copy operation below.
			<-ctx.Done()
			w.CloseWithError(errors.EnsureStack(ctx.Err()))
		}()
		if _, err := io.Copy(output, r); err != nil {
			return errors.EnsureStack(err)
		}
		return nil
	})
	reportProgress(100, 100)
	return filepath.Join(dir, "describe.txt"), err
}

func (s *debugServer) collectLogs(ctx context.Context, dir string, pod, container string, reportProgress func(progress, total int)) (string, error) {
	logsFile := "logs.txt"
	reportProgress(0, 1)
	if err := collectDebugFileV2(dir, logsFile, func(w io.Writer) (retErr error) {
		defer log.Span(ctx, "collectLogs", zap.String("pod", pod), zap.String("container", container))(log.Errorp(&retErr))
		stream, err := s.env.GetKubeClient().CoreV1().Pods(s.env.Config().Namespace).GetLogs(pod, &v1.PodLogOptions{Container: container}).Stream(ctx)
		if err != nil {
			return errors.EnsureStack(err)
		}
		defer func() {
			if err := stream.Close(); retErr == nil {
				retErr = err
			}
		}()
		_, err = io.Copy(w, stream)
		return errors.EnsureStack(err)
	}); err != nil {
		return filepath.Join(dir, logsFile), err
	}
	reportProgress(1, 2)
	logsFile = "logs-previous.txt"
	if err := collectDebugFileV2(dir, logsFile, func(w io.Writer) (retErr error) {
		defer log.Span(ctx, "collectLogs.previous", zap.String("pod", pod), zap.String("container", container))(log.Errorp(&retErr))
		stream, err := s.env.GetKubeClient().CoreV1().Pods(s.env.Config().Namespace).GetLogs(pod, &v1.PodLogOptions{Container: container, Previous: true}).Stream(ctx)
		if err != nil {
			return errors.EnsureStack(err)
		}
		defer func() {
			if err := stream.Close(); retErr == nil {
				retErr = err
			}
		}()
		_, err = io.Copy(w, stream)
		return errors.EnsureStack(err)
	}); err != nil {
		return logsFile, err
	}
	reportProgress(1, 1)
	// TODO: this isn't clean since it ignores the logs.txt file
	return logsFile, nil
}

func (s *debugServer) hasLoki() bool {
	_, err := s.env.GetLokiClient()
	return err == nil
}

func (s *debugServer) collectLogsLoki(ctx context.Context, dir, pod, container string, reportProgress func(progress, total int)) (string, error) {
	logsPath := filepath.Join(dir, "logs-loki.txt")
	if !s.hasLoki() {
		return "", nil
	}
	err := collectDebugFileV2(dir, "logs-loki.txt", func(w io.Writer) (retErr error) {
		defer log.Span(ctx, "collectLogsLoki", zap.String("pod", pod), zap.String("container", container))(log.Errorp(&retErr))
		queryStr := `{pod="` + pod
		if container != "" {
			queryStr += `", container="` + container
		}
		queryStr += `"}`
		logs, err := s.queryLoki(ctx, queryStr)
		if err != nil {
			return errors.EnsureStack(err)
		}
		progTotal := 100
		progIncr := len(logs)/progTotal + 1
		var cursor *loki.LabelSet
		for i, entry := range logs {
			// Print the stream labels in %v format whenever they are different from the
			// previous line.  The pointer comparison is a fast path to avoid
			// reflect.DeepEqual when both log lines are from the same chunk of logs
			// returned by Loki.
			if cursor != entry.Labels && !reflect.DeepEqual(cursor, entry.Labels) {
				if _, err := fmt.Fprintf(w, "%v\n", entry.Labels); err != nil {
					return errors.EnsureStack(err)
				}
			}
			cursor = entry.Labels
			// Then the line itself.
			if _, err := fmt.Fprintf(w, "%s\n", entry.Entry.Line); err != nil {
				return errors.EnsureStack(err)
			}
			if i%progIncr == 0 {
				reportProgress(i/progIncr, progTotal)
			}
		}
		reportProgress(100, 100)
		return nil
	})
	return logsPath, err
}

// collectDump is done on the pachd side, AND by the workers by them.
func collectDump(ctx context.Context, dir string, reportProgress func(progress, total int)) (_ string, retErr error) {
	defer log.Span(ctx, "collectDump", zap.String("path", dir))(log.Errorp(&retErr))
	reportProgress(0, 100)
	// Collect go info.
	if err := collectGoInfo(dir); err != nil {
		return dir, err
	}
	// Goroutine profile.
	if err := collectProfile(ctx, dir, &debug.Profile{Name: "goroutine"}); err != nil {
		return dir, err
	}
	// Heap profile.
	if err := collectProfile(ctx, dir, &debug.Profile{Name: "heap"}); err != nil {
		return dir, err
	}
	reportProgress(100, 100)
	return dir, nil
}

func (s *debugServer) collectPipelineV2(ctx context.Context, dir string, p *pps.Pipeline, limit int64) (retErr error) {
	ctx, end := log.SpanContext(ctx, "collectPipelineDump", zap.Stringer("pipeline", p))
	defer end(log.Errorp(&retErr))
	var errs error
	if err := collectDebugFileV2(dir, "spec.json", func(w io.Writer) error {
		fullPipelineInfos, err := s.env.GetPachClient(ctx).ListProjectPipelineHistory(p.Project.Name, p.Name, -1, true)
		if err != nil {
			return err
		}
		var pipelineErrs error
		for _, fullPipelineInfo := range fullPipelineInfos {
			if err := s.marshaller.Marshal(w, fullPipelineInfo); err != nil {
				multierr.AppendInto(&pipelineErrs, errors.Wrapf(err, "marshalFullPipelineInfo(%s)", fullPipelineInfo.GetPipeline()))
			}
		}
		return nil
	}); err != nil {
		multierr.AppendInto(&errs, errors.Wrap(err, "listProjectPipelineHistory"))
	}
	if err := s.collectCommits(ctx, s.env.GetPachClient(ctx), dir, client.NewProjectRepo(p.Project.GetName(), p.Name), limit); err != nil {
		multierr.AppendInto(&errs, errors.Wrap(err, "collectCommits"))
	}
	if err := s.collectJobs(s.env.GetPachClient(ctx), dir, p, limit); err != nil {
		multierr.AppendInto(&errs, errors.Wrap(err, "collectJobs"))
	}
	if s.hasLoki() {
		if err := s.forEachWorkerLoki(ctx, p, func(pod string) (retErr error) {
			ctx, end := log.SpanContext(ctx, "forEachWorkerLoki.worker", zap.String("pod", pod))
			defer end(log.Errorp(&retErr))
			workerDir := filepath.Join(dir, podPrefix, pod)
			userDir := filepath.Join(workerDir, client.PPSWorkerUserContainerName)
			if _, err := s.collectLogsLoki(ctx, userDir, pod, client.PPSWorkerUserContainerName, func(progress, total int) {}); err != nil {
				return err
			}
			sidecarDir := filepath.Join(workerDir, client.PPSWorkerSidecarContainerName)
			_, err := s.collectLogsLoki(ctx, sidecarDir, pod, client.PPSWorkerSidecarContainerName, func(progress, total int) {})
			return err
		}); err != nil {
			return writeErrorFileV2(filepath.Join(dir, "loki"), err)
		}
	}
	return errs
}

func (s *debugServer) forEachWorkerLoki(ctx context.Context, p *pps.Pipeline, cb func(string) error) error {
	pods, err := s.getWorkerPodsLoki(ctx, p)
	if err != nil {
		return err
	}
	var errs error
	for pod := range pods {
		if err := cb(pod); err != nil {
			multierr.AppendInto(&errs, errors.Wrapf(err, "forEachWorkersLoki(%s)", pod))
		}
	}
	return errs
}

// quoteLogQLStreamSelector returns a string quoted as a LogQL stream selector.
// The rules for quoting a LogQL stream selector are documented at
// https://grafana.com/docs/loki/latest/logql/log_queries/#log-stream-selector
// to be the same as for a Prometheus string literal, as documented at
// https://prometheus.io/docs/prometheus/latest/querying/basics/#string-literals.
// This happens to be the same as a Go string, with single or double quotes or
// backticks allowed as enclosing characters.
func quoteLogQLStreamSelector(s string) string {
	return strconv.Quote(s)
}

func (s *debugServer) getWorkerPodsLoki(ctx context.Context, p *pps.Pipeline) (map[string]struct{}, error) {
	// This function uses the log querying API and not the label querying API, to bound the the
	// number of workers for a pipeline that we return.  We'll get 30,000 of the most recent
	// logs for each pipeline, and return the names of the workers that contributed to those
	// logs for further inspection.  The alternative would be to get every worker that existed
	// in some time interval, but that results in too much data to inspect.

	queryStr := fmt.Sprintf(`{pipelineProject=%s, pipelineName=%s}`, quoteLogQLStreamSelector(p.Project.Name), quoteLogQLStreamSelector(p.Name))
	pods := make(map[string]struct{})
	logs, err := s.queryLoki(ctx, queryStr)
	if err != nil {
		return nil, err
	}

	for _, l := range logs {
		if l.Labels == nil {
			continue
		}
		labels := *l.Labels
		pod, ok := labels["pod"]
		if !ok {
			log.Debug(ctx, "pod label missing from loki labelset", zap.Any("log", l))
			continue
		}
		pods[pod] = struct{}{}
	}
	return pods, nil
}

const (
	// maxLogs used to be 5,000 but a comment said this was too few logs, so now it's 30,000.
	// If it still seems too small, bump it up again.
	maxLogs = 30000
	// 5,000 is the maximum value of "limit" in queries that Loki seems to accept.
	// We set serverMaxLogs below the actual limit because of a bug in Loki where it hangs when you request too much data from it.
	serverMaxLogs = 1000
)

type lokiLog struct {
	Labels *loki.LabelSet
	Entry  *loki.Entry
}

func (s *debugServer) queryLoki(ctx context.Context, queryStr string) (retResult []lokiLog, retErr error) {
	ctx, finishSpan := log.SpanContext(ctx, "queryLoki", zap.String("queryStr", queryStr))
	defer func() {
		finishSpan(zap.Error(retErr), zap.Int("logs", len(retResult)))
	}()

	sortLogs := func(logs []lokiLog) {
		sort.Slice(logs, func(i, j int) bool {
			return logs[i].Entry.Timestamp.Before(logs[j].Entry.Timestamp)
		})
	}

	c, err := s.env.GetLokiClient()
	if err != nil {
		return nil, errors.EnsureStack(errors.Errorf("get loki client: %v", err))
	}
	// We used to just stream the output, but that results in logs from different streams
	// (stdout/stderr) being randomly interspersed with each other, which is hard to follow when
	// written in text form.  We also used "FORWARD" and basically got the first 30,000 logs
	// this month, instead of the most recent 30,000 logs.  To use "BACKWARD" to get chunks of
	// logs starting with the most recent (whose start time we can't know without asking), we
	// have to either output logs in reverse order, or collect them all and then reverse.  Thus,
	// we just buffer.  30k logs * (200 bytes each + 16 bytes of pointers per line + 24 bytes of
	// time.Time) = 6.8MiB.  That seems totally reasonable to keep in memory, with the benefit
	// of providing lines in chronological order even when the app uses both stdout and stderr.
	var result []lokiLog

	var end time.Time
	start := time.Now().Add(-30 * 24 * time.Hour) // 30 days.  (Loki maximum range is 30 days + 1 hour.)

	for numLogs := 0; (end.IsZero() || start.Before(end)) && numLogs < maxLogs; {
		// Loki requests can hang if the size of the log lines is too big
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		resp, err := c.QueryRange(ctx, queryStr, serverMaxLogs, start, end, "BACKWARD", 0 /* step */, 0 /* interval */, true /* quiet */)
		cancel()
		if err != nil {
			// Note: the error from QueryRange has a stack.
			if errors.Is(err, context.DeadlineExceeded) {
				log.Debug(ctx, "loki query range timed out", zap.Int("serverMaxLogs", serverMaxLogs), zap.Time("start", start), zap.Time("end", end), zap.Int("logs", len(result)))
				sortLogs(result)
				return result, nil
			}
			return nil, errors.Errorf("query range (query=%v, limit=%v, start=%v, end=%v): %+v", queryStr, serverMaxLogs, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
		}

		streams, ok := resp.Data.Result.(loki.Streams)
		if !ok {
			return nil, errors.Errorf("resp.Data.Result must be of type loki.Streams")
		}

		var readThisIteration int
		for _, s := range streams {
			stream := s // Alias for pointer later.
			for _, e := range stream.Entries {
				entry := e
				numLogs++
				readThisIteration++
				if end.IsZero() || entry.Timestamp.Before(end) {
					// Because end is an "exclusive" range, if we read any logs
					// at all we are guaranteed to not read them in the next
					// iteration.  (If all 5000 logs are at the same time, we
					// still advance past them on the next iteration, which is
					// the best we can do under those circumstances.)
					end = entry.Timestamp
				}
				result = append(result, lokiLog{
					Labels: &stream.Labels,
					Entry:  &entry,
				})
			}
		}
		if readThisIteration == 0 {
			break
		}
		cancel()
	}
	sortLogs(result)
	return result, nil
}

func (s *debugServer) collectJobs(pachClient *client.APIClient, dir string, pipeline *pps.Pipeline, limit int64, prefix ...string) error {
	download := chart.ContinuousSeries{
		Name: "download",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(255),
		},
	}
	process := chart.ContinuousSeries{
		Name: "process",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(1).WithAlpha(255),
		},
	}
	upload := chart.ContinuousSeries{
		Name: "upload",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(2).WithAlpha(255),
		},
	}
	if err := collectDebugFileV2(dir, "jobs.json", func(w io.Writer) error {
		// TODO: The limiting should eventually be a feature of list job.
		var count int64
		return pachClient.ListProjectJobF(pipeline.Project.Name, pipeline.Name, nil, -1, false, func(ji *pps.JobInfo) error {
			if count >= limit {
				return errutil.ErrBreak
			}
			count++
			if ji.Stats.DownloadTime != nil && ji.Stats.ProcessTime != nil && ji.Stats.UploadTime != nil {
				downloadDuration, err := types.DurationFromProto(ji.Stats.DownloadTime)
				if err != nil {
					return errors.EnsureStack(err)
				}
				processDuration, err := types.DurationFromProto(ji.Stats.ProcessTime)
				if err != nil {
					return errors.EnsureStack(err)
				}
				uploadDuration, err := types.DurationFromProto(ji.Stats.UploadTime)
				if err != nil {
					return errors.EnsureStack(err)
				}
				download.XValues = append(download.XValues, float64(len(download.XValues)+1))
				download.YValues = append(download.YValues, float64(downloadDuration))
				process.XValues = append(process.XValues, float64(len(process.XValues)+1))
				process.YValues = append(process.YValues, float64(processDuration))
				upload.XValues = append(upload.XValues, float64(len(upload.XValues)+1))
				upload.YValues = append(upload.YValues, float64(uploadDuration))
			}
			return errors.EnsureStack(s.marshaller.Marshal(w, ji))
		})
	}); err != nil {
		return err
	}
	// Reverse the x values since we collect them from newest to oldest.
	// TODO: It would probably be better to support listing jobs in reverse order.
	reverseContinuousSeries(download, process, upload)
	return collectGraph(dir, "jobs-chart.png", "number of jobs", []chart.Series{download, process, upload}, prefix...)
}

// collectWorkerDump is something to do for each worker, on the pachd side.
func (s *debugServer) collectWorkerDump(ctx context.Context, dir string, pod *v1.Pod, prefix ...string) (retErr error) {
	defer log.Span(ctx, "collectWorkerDump")(log.Errorp(&retErr))
	// Collect the worker describe output.
	if _, err := s.collectDescribe(ctx, dir, pod.Name, func(_, _ int) {}); err != nil {
		return err
	}
	// Collect the worker user and storage container logs.
	userDir := filepath.Join(dir, client.PPSWorkerUserContainerName)
	sidecarDir := filepath.Join(dir, client.PPSWorkerSidecarContainerName)
	var errs error
	if _, err := s.collectLogs(ctx, userDir, pod.Name, client.PPSWorkerUserContainerName, func(_, _ int) {}); err != nil {
		multierr.AppendInto(&errs, errors.Wrap(err, "userContainerLogs"))
	}
	if _, err := s.collectLogs(ctx, sidecarDir, pod.Name, client.PPSWorkerSidecarContainerName, func(_, _ int) {}); err != nil {
		multierr.AppendInto(&errs, errors.Wrap(err, "storageContainerLogs"))
	}
	return errs
}

func redirectDump(ctx context.Context, c debug.DebugClient, filter *debug.Filter) (_ io.Reader, retErr error) {
	defer log.Span(ctx, "redirectDump")(log.Errorp(&retErr))
	dumpC, err := c.Dump(ctx, &debug.DumpRequest{Filter: filter})
	if err != nil {
		return nil, errors.EnsureStack(err)
	}
	return grpcutil.NewStreamingBytesReader(dumpC, nil), nil
}

func handleHelmSecretV2(ctx context.Context, path string, secret v1.Secret) error {
	var errs error
	name := secret.Name
	if got, want := string(secret.Type), "helm.sh/release.v1"; got != want {
		return errors.Errorf("helm-owned secret of unknown version; got %v want %v", got, want)
	}
	if secret.Data == nil {
		log.Info(ctx, "skipping helm secret with no data", zap.String("secretName", name))
		return nil
	}
	releaseData, ok := secret.Data["release"]
	if !ok {
		log.Info(ctx, "secret data doesn't have a release key", zap.String("secretName", name), zap.Strings("keys", maps.Keys(secret.Data)))
		return nil
	}
	// There are a few labels that helm adds that are interesting (the "release
	// status").  Grab those and write to metadata.json.
	if err := collectDebugFileV2(path, "metadata.json", func(w io.Writer) error {
		text, err := json.Marshal(secret.ObjectMeta)
		if err != nil {
			return errors.Wrap(err, "marshal ObjectMeta")
		}
		if _, err := w.Write(text); err != nil {
			return errors.Wrap(err, "write ObjectMeta json")
		}
		return nil
	}); err != nil {
		// We can still try to get the release JSON if the metadata doesn't marshal
		// or write cleanly.
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: collect metadata", name))
	}

	// Get the text of the release and write it to release.json.
	var releaseJSON []byte
	if err := collectDebugFileV2(path, "release.json", func(w io.Writer) error {
		// The helm release data is base64-encoded gzipped JSON-marshalled protobuf.
		// The base64 encoding is IN ADDITION to the base64 encoding that k8s does
		// when serving secrets through the API; client-go removes that for us.
		var err error
		releaseJSON, err = base64.StdEncoding.DecodeString(string(releaseData))
		if err != nil {
			releaseJSON = nil
			return errors.Wrap(err, "decode base64")
		}

		// Older versions of helm do not compress the data; this check is for the
		// gzip header; if deteced, decompress.
		if len(releaseJSON) > 3 && bytes.Equal(releaseJSON[0:3], []byte{0x1f, 0x8b, 0x08}) {
			gz, err := gzip.NewReader(bytes.NewReader(releaseJSON))
			if err != nil {
				return errors.Wrap(err, "create gzip reader")
			}
			releaseJSON, err = io.ReadAll(gz)
			if err != nil {
				return errors.Wrap(err, "decompress")
			}
			if err := gz.Close(); err != nil {
				return errors.Wrap(err, "close gzip reader")
			}
		}

		// Write out the raw release JSON, so fields we don't pick out in the next
		// phase can be parsed out by the person analyzing the dump if necessary.
		if _, err := w.Write(releaseJSON); err != nil {
			return errors.Wrap(err, "write release")
		}
		return nil
	}); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: collect release", name))
		// The next steps need the release JSON, so we have to give up if any of
		// this failed.  Technically if the write fails, we could continue, but it
		// doesn't seem worth the effort because the next writes are also likely to
		// fail.
		return errs
	}

	// Unmarshal the release JSON, and start writing files for each interesting part.
	var release struct {
		// Config is the merged values.yaml that aren't in the chart.
		Config map[string]any `json:"config"`
		// Manifest is the rendered manifest that was applied to the cluster.
		Manifest string `json:"manifest"`
	}
	if err := json.Unmarshal(releaseJSON, &release); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: unmarshal release json", name))
		return errs
	}

	// Write the manifest YAML.
	if err := collectDebugFileV2(path, "manifest.yaml", func(w io.Writer) error {
		// Helm adds a newline at the end of the manifest.
		if _, err := fmt.Fprintf(w, "%s", release.Manifest); err != nil {
			return errors.Wrap(err, "print manifest")
		}
		return nil
	}); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: write manifest.yaml", name))
		// We can try the next step if this fails.
	}

	// Write values.yaml.
	if err := collectDebugFileV2(path, "values.yaml", func(w io.Writer) error {
		b, err := yaml.Marshal(release.Config)
		if err != nil {
			return errors.Wrap(err, "marshal values to yaml")
		}
		if _, err := w.Write(b); err != nil {
			return errors.Wrap(err, "print values")
		}
		return nil
	}); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: write values.yaml", name))
	}
	return errs
}

// handleHelmSecret decodes a helm release secret into its various components.  A returned error
// should be written into error.txt.
func handleHelmSecret(ctx context.Context, tw *tar.Writer, secret v1.Secret) error {
	var errs error
	name := secret.Name
	if got, want := string(secret.Type), "helm.sh/release.v1"; got != want {
		return errors.Errorf("helm-owned secret of unknown version; got %v want %v", got, want)
	}
	if secret.Data == nil {
		log.Info(ctx, "skipping helm secret with no data", zap.String("secretName", name))
		return nil
	}
	releaseData, ok := secret.Data["release"]
	if !ok {
		log.Info(ctx, "secret data doesn't have a release key", zap.String("secretName", name), zap.Strings("keys", maps.Keys(secret.Data)))
		return nil
	}

	// There are a few labels that helm adds that are interesting (the "release
	// status").  Grab those and write to metadata.json.
	if err := collectDebugFile(tw, filepath.Join("helm", name, "metadata"), "json", func(w io.Writer) error {
		text, err := json.Marshal(secret.ObjectMeta)
		if err != nil {
			return errors.Wrap(err, "marshal ObjectMeta")
		}
		if _, err := w.Write(text); err != nil {
			return errors.Wrap(err, "write ObjectMeta json")
		}
		return nil
	}); err != nil {
		// We can still try to get the release JSON if the metadata doesn't marshal
		// or write cleanly.
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: collect metadata", name))
	}

	// Get the text of the release and write it to release.json.
	var releaseJSON []byte
	if err := collectDebugFile(tw, filepath.Join("helm", name, "release"), "json", func(w io.Writer) error {
		// The helm release data is base64-encoded gzipped JSON-marshalled protobuf.
		// The base64 encoding is IN ADDITION to the base64 encoding that k8s does
		// when serving secrets through the API; client-go removes that for us.
		var err error
		releaseJSON, err = base64.StdEncoding.DecodeString(string(releaseData))
		if err != nil {
			releaseJSON = nil
			return errors.Wrap(err, "decode base64")
		}

		// Older versions of helm do not compress the data; this check is for the
		// gzip header; if deteced, decompress.
		if len(releaseJSON) > 3 && bytes.Equal(releaseJSON[0:3], []byte{0x1f, 0x8b, 0x08}) {
			gz, err := gzip.NewReader(bytes.NewReader(releaseJSON))
			if err != nil {
				return errors.Wrap(err, "create gzip reader")
			}
			releaseJSON, err = io.ReadAll(gz)
			if err != nil {
				return errors.Wrap(err, "decompress")
			}
			if err := gz.Close(); err != nil {
				return errors.Wrap(err, "close gzip reader")
			}
		}

		// Write out the raw release JSON, so fields we don't pick out in the next
		// phase can be parsed out by the person analyzing the dump if necessary.
		if _, err := w.Write(releaseJSON); err != nil {
			return errors.Wrap(err, "write release")
		}
		return nil
	}); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: collect release", name))
		// The next steps need the release JSON, so we have to give up if any of
		// this failed.  Technically if the write fails, we could continue, but it
		// doesn't seem worth the effort because the next writes are also likely to
		// fail.
		return errs
	}

	// Unmarshal the release JSON, and start writing files for each interesting part.
	var release struct {
		// Config is the merged values.yaml that aren't in the chart.
		Config map[string]any `json:"config"`
		// Manifest is the rendered manifest that was applied to the cluster.
		Manifest string `json:"manifest"`
	}
	if err := json.Unmarshal(releaseJSON, &release); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: unmarshal release json", name))
		return errs
	}

	// Write the manifest YAML.
	if err := collectDebugFile(tw, filepath.Join("helm", name, "manifest"), "yaml", func(w io.Writer) error {
		// Helm adds a newline at the end of the manifest.
		if _, err := fmt.Fprintf(w, "%s", release.Manifest); err != nil {
			return errors.Wrap(err, "print manifest")
		}
		return nil
	}); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: write manifest.yaml", name))
		// We can try the next step if this fails.
	}

	// Write values.yaml.
	if err := collectDebugFile(tw, filepath.Join("helm", name, "values"), "yaml", func(w io.Writer) error {
		b, err := yaml.Marshal(release.Config)
		if err != nil {
			return errors.Wrap(err, "marshal values to yaml")
		}
		if _, err := w.Write(b); err != nil {
			return errors.Wrap(err, "print values")
		}
		return nil
	}); err != nil {
		multierr.AppendInto(&errs, errors.Wrapf(err, "%v: write values.yaml", name))
	}
	return errs
}

func (s *debugServer) helmReleases(ctx context.Context, tw *tar.Writer) (retErr error) {
	ctx, end := log.SpanContext(ctx, "collectHelmReleases")
	defer end(log.Errorp(&retErr))
	// Helm stores release data in secrets by default.  Users can override this by exporting
	// HELM_DRIVER and using something else, in which case, we won't get anything here.
	// https://helm.sh/docs/topics/advanced/#storage-backends
	secrets, err := s.env.GetKubeClient().CoreV1().Secrets(s.env.Config().Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "owner=helm",
	})
	if err != nil {
		return errors.EnsureStack(err)
	}
	var writeErrs error
	for _, secret := range secrets.Items {
		if err := handleHelmSecret(ctx, tw, secret); err != nil {
			if err := writeErrorFile(tw, err, filepath.Join("helm", secret.Name)); err != nil {
				multierr.AppendInto(&writeErrs, errors.Wrapf(err, "%v: write error.txt", secret.Name))
				continue
			}
		}
	}
	return writeErrs
}
