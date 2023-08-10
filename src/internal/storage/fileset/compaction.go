package fileset

import (
	"context"
	"math"
	"time"

	units "github.com/docker/go-units"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
	"github.com/pachyderm/pachyderm/v2/src/internal/errutil"
	"github.com/pachyderm/pachyderm/v2/src/internal/log"
	"github.com/pachyderm/pachyderm/v2/src/internal/pctx"
	"github.com/pachyderm/pachyderm/v2/src/internal/storage/fileset/index"
	"github.com/pachyderm/pachyderm/v2/src/internal/uuid"
)

// IsCompacted returns true if the file sets are already in a compacted form.
func (s *Storage) IsCompacted(ctx context.Context, h Handle) (bool, error) {
	var prev *Primitive
	compacted := true
	if err := s.Flatten(ctx, []Handle{h}, func(fsId ID) error {
		curr, err := s.getPrimitive(ctx, fsId)
		if err != nil {
			return err
		}
		if prev != nil && !s.isCompactedPair(prev, curr) {
			compacted = false
			return errutil.ErrBreak
		}
		prev = curr
		return nil
	}); err != nil {
		return false, err
	}
	return compacted, nil
}

func (s *Storage) isCompactedPair(left, right *Primitive) bool {
	return compactionScore(left) >= compactionScore(right)*s.compactionConfig.LevelFactor
}

func (s *Storage) isCompacted(prims []*Primitive) bool {
	return s.indexOfCompacted(prims) == len(prims)
}

// indexOfCompacted returns the last value of i for which the "compacted relationship" is maintained for all layers[:i+1].
// The "compacted relationship" is defined as leftScore >= (rightScore * factor).
// If there is an element at i+1, it will be the first element which does not satisfy the compacted relationship with i.
func (s *Storage) indexOfCompacted(prims []*Primitive) int {
	for i := 0; i < len(prims)-1; i++ {
		if !s.isCompactedPair(prims[i], prims[i+1]) {
			return i
		}
	}
	return len(prims)
}

// indexOfCompactedOptimized is like indexOfCompacted but runs an optimization to reduce extra compaction operations.
func (s *Storage) indexOfCompactedOptimized(prims []*Primitive) int {
	for i := 0; i < len(prims)-1; i++ {
		if !s.isCompactedPair(prims[i], prims[i+1]) {
			var score int64
			for j := i; j < len(prims); j++ {
				score += compactionScore(prims[j])
			}
			for ; i > 0; i-- {
				if compactionScore(prims[i-1]) >= score*s.compactionConfig.LevelFactor {
					return i
				}
				score += compactionScore(prims[i-1])
			}
			return i
		}
	}
	return len(prims)
}

// compactionScore computes a score for a primitive file set that can be used for making decisions about compaction.
// A higher score means that more work would be involved when including the associated primitive file set in a compaction.
func compactionScore(prim *Primitive) int64 {
	// TODO: Add to prevent full compaction when migrating?
	//	if prim.SizeBytes > 0 {
	//		return prim.SizeBytes
	//	}
	var score int64
	if prim.Additive != nil {
		score += prim.Additive.NumFiles * units.KB
		score += prim.Additive.SizeBytes
	}
	if prim.Deletive != nil {
		score += prim.Deletive.NumFiles * units.KB
	}
	return score
}

// Compact compacts the contents of ids into a new file set with the specified ttl and returns the ID.
// Compact always returns the ID of a primitive file set.
// Compact does not renew ids.
// It is the responsibility of the caller to renew ids.  In some cases they may be permanent and not require renewal.
func (s *Storage) Compact(ctx context.Context, hs []Handle, ttl time.Duration, opts ...index.Option) (*Handle, error) {
	w := s.newWriter(ctx, WithTTL(ttl))
	fs, err := s.Open(ctx, hs)
	if err != nil {
		return nil, err
	}
	// TODO: Consider adding prefetching here.
	if err := CopyDeletedFiles(ctx, w, fs, opts...); err != nil {
		return nil, err
	}
	if err := CopyFiles(ctx, w, fs, opts...); err != nil {
		return nil, err
	}
	return w.Close()
}

// CompactCallback is the standard callback signature for a compaction operation.
type CompactCallback func(context.Context, []Handle, time.Duration) (*Handle, error)

// CompactLevelBased performs a level-based compaction on the passed in file sets.
func (s *Storage) CompactLevelBased(ctx context.Context, hs []Handle, maxFanIn int, ttl time.Duration, compact CompactCallback) (*Handle, error) {
	ids, err := s.FlattenAll(ctx, hs)
	if err != nil {
		return nil, err
	}
	prims, err := s.getPrimitives(ctx, ids)
	if err != nil {
		return nil, err
	}
	if s.isCompacted(prims) {
		return s.Compose(ctx, hs, ttl)
	}
	var h *Handle
	if err := s.WithRenewer(ctx, ttl, func(ctx context.Context, renewer *Renewer) error {
		i := s.indexOfCompactedOptimized(prims)
		if err := log.LogStep(ctx, "compactLevels", func(ctx context.Context) error {
			h, err = s.compactLevels(ctx, hs[i:], maxFanIn, ttl, compact)
			if err != nil {
				return err
			}
			return renewer.Add(ctx, *h)
		}, zap.Int("indexOfCompactedOptimized", i), zap.Int("ids", len(hs))); err != nil {
			return err
		}
		h, err = s.CompactLevelBased(ctx, append(hs[:i], *h), maxFanIn, ttl, compact)
		return err
	}); err != nil {
		return nil, err
	}
	return h, nil
}

// compactLevels compacts a list of levels.
// The compaction happens in steps where each step includes file sets of similar compaction score.
// For each step, ranges of file sets in the list that are of length maxFanIn and contain file sets that are
// less than or equal to the step score are compacted.
// For each step, this process is repeated until there are no more ranges eligible for compaction in the step.
// The file sets being compacted must be contiguous because the file operation order matters.
// This algorithm ensures that we compact file sets with lower scores together first before compacting higher score file sets.
func (s *Storage) compactLevels(ctx context.Context, hs []Handle, maxFanIn int, ttl time.Duration, compact CompactCallback) (*Handle, error) {
	var h *Handle
	if err := s.WithRenewer(ctx, ttl, func(ctx context.Context, renewer *Renewer) error {
		for step := 0; len(hs) > maxFanIn; step++ {
			if err := log.LogStep(ctx, "compactLevels.step", func(ctx context.Context) error {
				stepScore := s.stepScore(step)
				var emptyStep bool
				for !emptyStep {
					emptyStep = true
					nextIds := make([]Handle, 0, len(hs))
					var compactHs []Handle
					eg, ctx := errgroup.WithContext(ctx)
					for _, h := range hs {
						prim, err := s.getPrimitive(ctx, h.id)
						if err != nil {
							return err
						}
						compactHs = append(compactHs, h)
						if compactionScore(prim) > stepScore {
							nextIds = append(nextIds, compactHs...)
							compactHs = nil
							continue
						}
						if len(compactHs) == maxFanIn {
							emptyStep = false
							hs := compactHs
							i := len(nextIds)
							nextIds = append(nextIds, Handle{})
							eg.Go(func() error {
								return log.LogStep(ctx, "compactBatch", func(ctx context.Context) error {
									h, err := compact(ctx, hs, ttl)
									if err != nil {
										return err
									}
									if err := renewer.Add(ctx, *h); err != nil {
										return err
									}
									nextIds[i] = *h
									return nil
								}, zap.String("batch", uuid.NewWithoutDashes()))
							})
							compactHs = nil
						}
					}
					nextIds = append(nextIds, compactHs...)
					if err := eg.Wait(); err != nil {
						return errors.EnsureStack(err)
					}
					hs = nextIds
				}
				return nil
			}, zap.Int("step", step)); err != nil {
				return err
			}
		}
		if len(hs) == 1 {
			h = &hs[0]
			return nil
		}
		var err error
		h, err = compact(pctx.Child(ctx, "compact", pctx.WithFields(zap.String("batch", uuid.NewWithoutDashes()))), hs, ttl)
		return err
	}); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Storage) stepScore(step int) int64 {
	return s.shardConfig.SizeBytes * int64(math.Pow(float64(s.compactionConfig.LevelFactor), float64(step)))
}
