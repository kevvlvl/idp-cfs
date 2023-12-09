package platform_git

import (
	"flag"
	"github.com/google/go-github/v56/github"
)

type Git interface {
	GetOrganization(organizationName string) (*Organization, error)
	GetRepository(repoName string) (*Repository, error)
	CreateRepository(repoName string, branch string) (*Repository, error)
}

type GitCode struct {
	githubCode   *GithubCode
	Repository   *Repository
	Organization *Organization
}

type GithubCode struct {
	GithubClient *github.Client
	githubUser   *github.User
}

type Organization struct {
	Name    *string
	Company *string
}

type Repository struct {
	Name         *string
	Organization *string
	Owner        *string
	URL          *string
}

const (
	// CodeGithub for the code repository of type Github (public/cloud)
	CodeGithub = "github"
	// CodeGitlab for the code repository of type Gitlab
	CodeGitlab = "gitlab"
	// CodeGitea for the code repository of type Gitea
	CodeGitea = "gitea"
)

var GithubPrivateRepository = flag.Bool("private", false, "Will created repo be private.")
var GithubDescription = flag.String("description", "", "Created by idp-cfs")
var GithubAutoInit = flag.Bool("auto-init", false, "Pass true to create an initial commit with empty README.")
