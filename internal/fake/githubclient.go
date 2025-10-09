package fake

import (
	"context"

	"github.com/hermesgen/hm"
)

// GithubClient is a fake implementation of hm.GitClient for testing.
type GithubClient struct {
	// Expected results
	CloneFn    func(ctx context.Context, repoURL, localPath string, auth hm.GitAuth, env []string) error
	CheckoutFn func(ctx context.Context, localRepoPath, branch string, create bool, env []string) error
	AddFn      func(ctx context.Context, localRepoPath, pathspec string, env []string) error
	CommitFn   func(ctx context.Context, localRepoPath string, commit hm.GitCommit, env []string) (string, error)
	PushFn     func(ctx context.Context, localRepoPath string, auth hm.GitAuth, remote, branch string, env []string) error
	StatusFn   func(ctx context.Context, localRepoPath string, env []string) (string, error)
	GitLogFn   func(ctx context.Context, localRepoPath string, args []string, env []string) (string, error)

	// Captured arguments
	CloneCalls []struct {
		Ctx                context.Context
		RepoURL, LocalPath string
		Auth               hm.GitAuth
		Env                []string
	}
	CheckoutCalls []struct {
		Ctx                   context.Context
		LocalRepoPath, Branch string
		Create                bool
		Env                   []string
	}
	AddCalls []struct {
		Ctx                     context.Context
		LocalRepoPath, Pathspec string
		Env                     []string
	}
	CommitCalls []struct {
		Ctx           context.Context
		LocalRepoPath string
		Commit        hm.GitCommit
		Env           []string
	}
	PushCalls []struct {
		Ctx                           context.Context
		LocalRepoPath, Remote, Branch string
		Auth                          hm.GitAuth
		Env                           []string
	}
	StatusCalls []struct {
		Ctx           context.Context
		LocalRepoPath string
		Env           []string
	}
	GitLogCalls []struct {
		Ctx           context.Context
		LocalRepoPath string
		Args          []string
		Env           []string
	}
}

// NewGithubClient creates a new fake GithubClient.
func NewGithubClient() *GithubClient {
	return &GithubClient{}
}

func (f *GithubClient) Clone(ctx context.Context, repoURL, localPath string, auth hm.GitAuth, env []string) error {
	f.CloneCalls = append(f.CloneCalls, struct {
		Ctx                context.Context
		RepoURL, LocalPath string
		Auth               hm.GitAuth
		Env                []string
	}{Ctx: ctx, RepoURL: repoURL, LocalPath: localPath, Auth: auth, Env: env})
	if f.CloneFn != nil {
		return f.CloneFn(ctx, repoURL, localPath, auth, env)
	}
	return nil // Default success
}

func (f *GithubClient) Checkout(ctx context.Context, localRepoPath, branch string, create bool, env []string) error {
	f.CheckoutCalls = append(f.CheckoutCalls, struct {
		Ctx                   context.Context
		LocalRepoPath, Branch string
		Create                bool
		Env                   []string
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Branch: branch, Create: create, Env: env})
	if f.CheckoutFn != nil {
		return f.CheckoutFn(ctx, localRepoPath, branch, create, env)
	}
	return nil
}

func (f *GithubClient) Add(ctx context.Context, localRepoPath, pathspec string, env []string) error {
	f.AddCalls = append(f.AddCalls, struct {
		Ctx                     context.Context
		LocalRepoPath, Pathspec string
		Env                     []string
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Pathspec: pathspec, Env: env})
	if f.AddFn != nil {
		return f.AddFn(ctx, localRepoPath, pathspec, env)
	}
	return nil
}

func (f *GithubClient) Commit(ctx context.Context, localRepoPath string, commit hm.GitCommit, env []string) (string, error) {
	f.CommitCalls = append(f.CommitCalls, struct {
		Ctx           context.Context
		LocalRepoPath string
		Commit        hm.GitCommit
		Env           []string
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Commit: commit, Env: env})
	if f.CommitFn != nil {
		return f.CommitFn(ctx, localRepoPath, commit, env)
	}
	return "fake-commit-hash", nil // Default hash
}

func (f *GithubClient) Push(ctx context.Context, localRepoPath string, auth hm.GitAuth, remote, branch string, env []string) error {
	f.PushCalls = append(f.PushCalls, struct {
		Ctx                           context.Context
		LocalRepoPath, Remote, Branch string
		Auth                          hm.GitAuth
		Env                           []string
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Auth: auth, Remote: remote, Branch: branch, Env: env})
	if f.PushFn != nil {
		return f.PushFn(ctx, localRepoPath, auth, remote, branch, env)
	}
	return nil
}

func (f *GithubClient) Status(ctx context.Context, localRepoPath string, env []string) (string, error) {
	f.StatusCalls = append(f.StatusCalls, struct {
		Ctx           context.Context
		LocalRepoPath string
		Env           []string
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Env: env})
	if f.StatusFn != nil {
		return f.StatusFn(ctx, localRepoPath, env)
	}
	return " M somefile.txt\n?? anotherfile.txt", nil // Default status
}

func (f *GithubClient) GitLog(ctx context.Context, localRepoPath string, args []string, env []string) (string, error) {
	f.GitLogCalls = append(f.GitLogCalls, struct {
		Ctx           context.Context
		LocalRepoPath string
		Args          []string
		Env           []string
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Args: args, Env: env})
	if f.GitLogFn != nil {
		return f.GitLogFn(ctx, localRepoPath, args, env)
	}
	return "fake git log", nil // Default log
}
