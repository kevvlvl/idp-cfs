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
	PushGoldenPath(url, pathDir, branch, gpWorkdir, codeWorkDir string, tag *string) error
}

type GithubApi struct {
	client      *github.Client
	user        *github.User
	ctx         context.Context
	repository  *github.Repository
	getRepoFunc func(ctx context.Context, c *github.Client, owner, repo string) (*github.Repository, *github.Response, error)
}

type GitlabApi struct {
	client  *gitlab.Client
	project *gitlab.Project
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
