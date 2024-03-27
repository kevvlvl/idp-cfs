package git_api

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/xanzy/go-gitlab"
)

type GitSource interface {
	ValidateNewCode(repoName string) error
	ValidateUpdateCode(repoName string) error
	ValidateGoldenPath(url, branch, workDir string) error
	CreateRepo(repoName string) error
	PushGoldenPath(url, pathDir, codeBranch, branch, gpWorkdir, codeWorkDir string, tag *string) error
	UpdateRepo(repoName string, newGpBranch string) error
	OpenPullRequest(repoName, refBranch, baseBranch string) error
}

type GithubApi struct {
	client         *github.Client
	user           *github.User
	ctx            context.Context
	repository     *github.Repository
	getRepoFunc    func(ctx context.Context, c *github.Client, owner, repo string) (*github.Repository, *github.Response, error)
	createRepoFunc func(ctx context.Context, c *github.Client, org string, repo *github.Repository) (*github.Repository, *github.Response, error)
	createFileFunc func(ctx context.Context, c *github.Client, owner, repo, path string, opts *github.RepositoryContentFileOptions) error
	createRefFunc  func(ctx context.Context, c *github.Client, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
	createPrFunc   func(ctx context.Context, c *github.Client, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
}

type GitlabApi struct {
	client            *gitlab.Client
	project           *gitlab.Project
	newClientFunc     func(token string, options ...gitlab.ClientOptionFunc) (*gitlab.Client, error)
	getProjectFunc    func(c *gitlab.Client, pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error)
	createProjectFunc func(c *gitlab.Client, opt *gitlab.CreateProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error)
}

type GitApiAuth struct {
	codeDefined bool
	codeUser    string
	codePass    string
	codeToken   string
	gpDefined   bool
	gpUser      string
	gpPass      string
	gpToken     string
}

const (
	ToolGithub = "GITHUB"
	ToolGitlab = "GITLAB"
)
