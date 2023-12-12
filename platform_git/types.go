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
	githubCode    *GithubCode
	Repository    *Repository
	Organization  *Organization
	CodeClonePath string
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
	// GitCommitAuthor is the author of the git commit
	GitCommitAuthor = "idp-cfs"
	// GitCommitAuthorEmail is the email of the git commit's author
	GitCommitAuthorEmail = "idp-cfs@kevvlvl.github.noreply.com"
	// CodeClonePath is the path where we git clone the code path to prepare pushing the golden path (when defined)
	CodeClonePath = "/tmp/idp-cfs-code"
)

var GithubPrivateRepository = flag.Bool("private", false, "Will created repo be private.")
var GithubDescription = flag.String("description", "", "Created by idp-cfs")
var GithubAutoInit = flag.Bool("auto-init", false, "Pass true to create an initial commit with empty README.")
