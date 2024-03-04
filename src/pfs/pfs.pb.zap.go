// Code generated by protoc-gen-zap (etc/proto/protoc-gen-zap). DO NOT EDIT.
//
// source: pfs/pfs.proto

package pfs

import (
	fmt "fmt"
	protoextensions "github.com/pachyderm/pachyderm/v2/src/protoextensions"
	zapcore "go.uber.org/zap/zapcore"
)

func (x *Repo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("name", x.Name)
	enc.AddString("type", x.Type)
	enc.AddObject("project", x.Project)
	return nil
}

func (x *RepoPicker) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("name", x.GetName())
	return nil
}

func (x *RepoPicker_RepoName) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("project", x.Project)
	enc.AddString("name", x.Name)
	enc.AddString("type", x.Type)
	return nil
}

func (x *Branch) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddString("name", x.Name)
	return nil
}

func (x *BranchPicker) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("name", x.GetName())
	return nil
}

func (x *BranchPicker_BranchName) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddString("name", x.Name)
	return nil
}

func (x *File) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddString("path", x.Path)
	enc.AddString("datum", x.Datum)
	return nil
}

func (x *RepoInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	protoextensions.AddTimestamp(enc, "created", x.Created)
	enc.AddInt64("size_bytes_upper_bound", x.SizeBytesUpperBound)
	enc.AddString("description", x.Description)
	branchesArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Branches {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("branches", zapcore.ArrayMarshalerFunc(branchesArrMarshaller))
	enc.AddObject("auth_info", x.AuthInfo)
	enc.AddObject("details", x.Details)
	return nil
}

func (x *RepoInfo_Details) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddInt64("size_bytes", x.SizeBytes)
	return nil
}

func (x *AuthInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	permissionsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Permissions {
			enc.AppendString(v.String())
		}
		return nil
	}
	enc.AddArray("permissions", zapcore.ArrayMarshalerFunc(permissionsArrMarshaller))
	rolesArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Roles {
			enc.AppendString(v)
		}
		return nil
	}
	enc.AddArray("roles", zapcore.ArrayMarshalerFunc(rolesArrMarshaller))
	return nil
}

func (x *BranchInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("branch", x.Branch)
	enc.AddObject("head", x.Head)
	provenanceArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Provenance {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("provenance", zapcore.ArrayMarshalerFunc(provenanceArrMarshaller))
	subvenanceArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Subvenance {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("subvenance", zapcore.ArrayMarshalerFunc(subvenanceArrMarshaller))
	direct_provenanceArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.DirectProvenance {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("direct_provenance", zapcore.ArrayMarshalerFunc(direct_provenanceArrMarshaller))
	enc.AddObject("trigger", x.Trigger)
	return nil
}

func (x *Trigger) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("branch", x.Branch)
	enc.AddBool("all", x.All)
	enc.AddString("rate_limit_spec", x.RateLimitSpec)
	enc.AddString("size", x.Size)
	enc.AddInt64("commits", x.Commits)
	enc.AddString("cron_spec", x.CronSpec)
	return nil
}

func (x *CommitOrigin) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("kind", x.Kind.String())
	return nil
}

func (x *Commit) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddString("id", x.Id)
	enc.AddObject("branch", x.Branch)
	return nil
}

func (x *CommitPicker) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("branch_head", x.GetBranchHead())
	enc.AddObject("id", x.GetId())
	enc.AddObject("ancestor", x.GetAncestor())
	enc.AddObject("branch_root", x.GetBranchRoot())
	return nil
}

func (x *CommitPicker_CommitByGlobalId) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddString("id", x.Id)
	return nil
}

func (x *CommitPicker_BranchRoot) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddUint32("offset", x.Offset)
	enc.AddObject("branch", x.Branch)
	return nil
}

func (x *CommitPicker_AncestorOf) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddUint32("offset", x.Offset)
	enc.AddObject("start", x.Start)
	return nil
}

func (x *CommitInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddObject("origin", x.Origin)
	enc.AddString("description", x.Description)
	enc.AddObject("parent_commit", x.ParentCommit)
	child_commitsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.ChildCommits {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("child_commits", zapcore.ArrayMarshalerFunc(child_commitsArrMarshaller))
	protoextensions.AddTimestamp(enc, "started", x.Started)
	protoextensions.AddTimestamp(enc, "finishing", x.Finishing)
	protoextensions.AddTimestamp(enc, "finished", x.Finished)
	direct_provenanceArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.DirectProvenance {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("direct_provenance", zapcore.ArrayMarshalerFunc(direct_provenanceArrMarshaller))
	enc.AddString("error", x.Error)
	enc.AddInt64("size_bytes_upper_bound", x.SizeBytesUpperBound)
	enc.AddObject("details", x.Details)
	return nil
}

func (x *CommitInfo_Details) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddInt64("size_bytes", x.SizeBytes)
	protoextensions.AddDuration(enc, "compacting_time", x.CompactingTime)
	protoextensions.AddDuration(enc, "validating_time", x.ValidatingTime)
	return nil
}

func (x *CommitSet) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("id", x.Id)
	return nil
}

func (x *CommitSetInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit_set", x.CommitSet)
	commitsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Commits {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("commits", zapcore.ArrayMarshalerFunc(commitsArrMarshaller))
	return nil
}

func (x *FileInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("file", x.File)
	enc.AddString("file_type", x.FileType.String())
	protoextensions.AddTimestamp(enc, "committed", x.Committed)
	enc.AddInt64("size_bytes", x.SizeBytes)
	protoextensions.AddBytes(enc, "hash", x.Hash)
	return nil
}

func (x *Project) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("name", x.Name)
	return nil
}

func (x *ProjectInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("project", x.Project)
	enc.AddString("description", x.Description)
	enc.AddObject("auth_info", x.AuthInfo)
	protoextensions.AddTimestamp(enc, "created_at", x.CreatedAt)
	return nil
}

func (x *ProjectPicker) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("name", x.GetName())
	return nil
}

func (x *CreateRepoRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddString("description", x.Description)
	enc.AddBool("update", x.Update)
	return nil
}

func (x *InspectRepoRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	return nil
}

func (x *ListRepoRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("type", x.Type)
	projectsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Projects {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("projects", zapcore.ArrayMarshalerFunc(projectsArrMarshaller))
	enc.AddObject("page", x.Page)
	return nil
}

func (x *RepoPage) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("order", x.Order.String())
	enc.AddInt64("page_size", x.PageSize)
	enc.AddInt64("page_index", x.PageIndex)
	return nil
}

func (x *DeleteRepoRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddBool("force", x.Force)
	return nil
}

func (x *DeleteReposRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	projectsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Projects {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("projects", zapcore.ArrayMarshalerFunc(projectsArrMarshaller))
	enc.AddBool("force", x.Force)
	enc.AddBool("all", x.All)
	return nil
}

func (x *DeleteRepoResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddBool("deleted", x.Deleted)
	return nil
}

func (x *DeleteReposResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	reposArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Repos {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("repos", zapcore.ArrayMarshalerFunc(reposArrMarshaller))
	return nil
}

func (x *StartCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("parent", x.Parent)
	enc.AddString("description", x.Description)
	enc.AddObject("branch", x.Branch)
	return nil
}

func (x *FinishCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddString("description", x.Description)
	enc.AddString("error", x.Error)
	enc.AddBool("force", x.Force)
	return nil
}

func (x *InspectCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddString("wait", x.Wait.String())
	return nil
}

func (x *ListCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddObject("from", x.From)
	enc.AddObject("to", x.To)
	enc.AddInt64("number", x.Number)
	enc.AddBool("reverse", x.Reverse)
	enc.AddBool("all", x.All)
	enc.AddString("origin_kind", x.OriginKind.String())
	protoextensions.AddTimestamp(enc, "started_time", x.StartedTime)
	return nil
}

func (x *InspectCommitSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit_set", x.CommitSet)
	enc.AddBool("wait", x.Wait)
	return nil
}

func (x *ListCommitSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("project", x.Project)
	return nil
}

func (x *SquashCommitSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit_set", x.CommitSet)
	return nil
}

func (x *DropCommitSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit_set", x.CommitSet)
	return nil
}

func (x *SubscribeCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddString("branch", x.Branch)
	enc.AddObject("from", x.From)
	enc.AddString("state", x.State.String())
	enc.AddBool("all", x.All)
	enc.AddString("origin_kind", x.OriginKind.String())
	return nil
}

func (x *ClearCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	return nil
}

func (x *SquashCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddBool("recursive", x.Recursive)
	return nil
}

func (x *SquashCommitResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	return nil
}

func (x *DropCommitRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddBool("recursive", x.Recursive)
	return nil
}

func (x *DropCommitResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	return nil
}

func (x *CreateBranchRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("head", x.Head)
	enc.AddObject("branch", x.Branch)
	provenanceArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Provenance {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("provenance", zapcore.ArrayMarshalerFunc(provenanceArrMarshaller))
	enc.AddObject("trigger", x.Trigger)
	enc.AddBool("new_commit_set", x.NewCommitSet)
	return nil
}

func (x *FindCommitsRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("start", x.Start)
	enc.AddString("file_path", x.FilePath)
	enc.AddUint32("limit", x.Limit)
	return nil
}

func (x *FindCommitsResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("found_commit", x.GetFoundCommit())
	enc.AddObject("last_searched_commit", x.GetLastSearchedCommit())
	enc.AddUint32("commits_searched", x.CommitsSearched)
	return nil
}

func (x *InspectBranchRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("branch", x.Branch)
	return nil
}

func (x *ListBranchRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("repo", x.Repo)
	enc.AddBool("reverse", x.Reverse)
	return nil
}

func (x *DeleteBranchRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("branch", x.Branch)
	enc.AddBool("force", x.Force)
	return nil
}

func (x *CreateProjectRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("project", x.Project)
	enc.AddString("description", x.Description)
	enc.AddBool("update", x.Update)
	return nil
}

func (x *InspectProjectRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("project", x.Project)
	return nil
}

func (x *InspectProjectV2Request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("project", x.Project)
	return nil
}

func (x *InspectProjectV2Response) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("info", x.Info)
	enc.AddString("defaults_json", x.DefaultsJson)
	return nil
}

func (x *ListProjectRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	return nil
}

func (x *DeleteProjectRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("project", x.Project)
	enc.AddBool("force", x.Force)
	return nil
}

func (x *AddFile) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("path", x.Path)
	enc.AddString("datum", x.Datum)
	protoextensions.AddBytesValue(enc, "raw", x.GetRaw())
	enc.AddObject("url", x.GetUrl())
	return nil
}

func (x *AddFile_URLSource) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("URL", x.URL)
	enc.AddBool("recursive", x.Recursive)
	enc.AddUint32("concurrency", x.Concurrency)
	return nil
}

func (x *DeleteFile) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("path", x.Path)
	enc.AddString("datum", x.Datum)
	return nil
}

func (x *CopyFile) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("dst", x.Dst)
	enc.AddString("datum", x.Datum)
	enc.AddObject("src", x.Src)
	enc.AddBool("append", x.Append)
	return nil
}

func (x *ModifyFileRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("set_commit", x.GetSetCommit())
	enc.AddObject("add_file", x.GetAddFile())
	enc.AddObject("delete_file", x.GetDeleteFile())
	enc.AddObject("copy_file", x.GetCopyFile())
	return nil
}

func (x *GetFileRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("file", x.File)
	enc.AddString("URL", x.URL)
	enc.AddInt64("offset", x.Offset)
	enc.AddObject("path_range", x.PathRange)
	return nil
}

func (x *InspectFileRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("file", x.File)
	return nil
}

func (x *ListFileRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("file", x.File)
	enc.AddObject("paginationMarker", x.PaginationMarker)
	enc.AddInt64("number", x.Number)
	enc.AddBool("reverse", x.Reverse)
	return nil
}

func (x *WalkFileRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("file", x.File)
	enc.AddObject("paginationMarker", x.PaginationMarker)
	enc.AddInt64("number", x.Number)
	enc.AddBool("reverse", x.Reverse)
	return nil
}

func (x *GlobFileRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddString("pattern", x.Pattern)
	enc.AddObject("path_range", x.PathRange)
	return nil
}

func (x *DiffFileRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("new_file", x.NewFile)
	enc.AddObject("old_file", x.OldFile)
	enc.AddBool("shallow", x.Shallow)
	return nil
}

func (x *DiffFileResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("new_file", x.NewFile)
	enc.AddObject("old_file", x.OldFile)
	return nil
}

func (x *FsckRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddBool("fix", x.Fix)
	enc.AddObject("zombie_target", x.GetZombieTarget())
	enc.AddBool("zombie_all", x.GetZombieAll())
	return nil
}

func (x *FsckResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("fix", x.Fix)
	enc.AddString("error", x.Error)
	return nil
}

func (x *CreateFileSetResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("file_set_id", x.FileSetId)
	return nil
}

func (x *GetFileSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddString("type", x.Type.String())
	return nil
}

func (x *AddFileSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddString("file_set_id", x.FileSetId)
	return nil
}

func (x *RenewFileSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("file_set_id", x.FileSetId)
	enc.AddInt64("ttl_seconds", x.TtlSeconds)
	return nil
}

func (x *ComposeFileSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	file_set_idsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.FileSetIds {
			enc.AppendString(v)
		}
		return nil
	}
	enc.AddArray("file_set_ids", zapcore.ArrayMarshalerFunc(file_set_idsArrMarshaller))
	enc.AddInt64("ttl_seconds", x.TtlSeconds)
	enc.AddBool("compact", x.Compact)
	return nil
}

func (x *ShardFileSetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("file_set_id", x.FileSetId)
	enc.AddInt64("num_files", x.NumFiles)
	enc.AddInt64("size_bytes", x.SizeBytes)
	return nil
}

func (x *PathRange) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("lower", x.Lower)
	enc.AddString("upper", x.Upper)
	return nil
}

func (x *ShardFileSetResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	shardsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Shards {
			enc.AppendObject(v)
		}
		return nil
	}
	enc.AddArray("shards", zapcore.ArrayMarshalerFunc(shardsArrMarshaller))
	return nil
}

func (x *CheckStorageRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddBool("read_chunk_data", x.ReadChunkData)
	protoextensions.AddBytes(enc, "chunk_begin", x.ChunkBegin)
	protoextensions.AddBytes(enc, "chunk_end", x.ChunkEnd)
	return nil
}

func (x *CheckStorageResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddInt64("chunk_object_count", x.ChunkObjectCount)
	return nil
}

func (x *PutCacheRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("key", x.Key)
	protoextensions.AddAny(enc, "value", x.Value)
	file_set_idsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.FileSetIds {
			enc.AppendString(v)
		}
		return nil
	}
	enc.AddArray("file_set_ids", zapcore.ArrayMarshalerFunc(file_set_idsArrMarshaller))
	enc.AddString("tag", x.Tag)
	return nil
}

func (x *GetCacheRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("key", x.Key)
	return nil
}

func (x *GetCacheResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	protoextensions.AddAny(enc, "value", x.Value)
	return nil
}

func (x *ClearCacheRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("tag_prefix", x.TagPrefix)
	return nil
}

func (x *ActivateAuthRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	return nil
}

func (x *ActivateAuthResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	return nil
}

func (x *ObjectStorageEgress) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("url", x.Url)
	return nil
}

func (x *SQLDatabaseEgress) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("url", x.Url)
	enc.AddObject("file_format", x.FileFormat)
	enc.AddObject("secret", x.Secret)
	return nil
}

func (x *SQLDatabaseEgress_FileFormat) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("type", x.Type.String())
	columnsArrMarshaller := func(enc zapcore.ArrayEncoder) error {
		for _, v := range x.Columns {
			enc.AppendString(v)
		}
		return nil
	}
	enc.AddArray("columns", zapcore.ArrayMarshalerFunc(columnsArrMarshaller))
	return nil
}

func (x *SQLDatabaseEgress_Secret) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddString("name", x.Name)
	enc.AddString("key", x.Key)
	return nil
}

func (x *EgressRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("commit", x.Commit)
	enc.AddObject("object_storage", x.GetObjectStorage())
	enc.AddObject("sql_database", x.GetSqlDatabase())
	return nil
}

func (x *EgressResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("object_storage", x.GetObjectStorage())
	enc.AddObject("sql_database", x.GetSqlDatabase())
	return nil
}

func (x *EgressResponse_ObjectStorageResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddInt64("bytes_written", x.BytesWritten)
	return nil
}

func (x *EgressResponse_SQLDatabaseResult) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if x == nil {
		return nil
	}
	enc.AddObject("rows_written", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for k, v := range x.RowsWritten {
			enc.AddInt64(fmt.Sprintf("%v", k), v)
		}
		return nil
	}))
	return nil
}
