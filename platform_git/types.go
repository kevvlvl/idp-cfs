package platform_git

import "github.com/google/go-github/v56/github"

type Git interface {
	GetOrganization(organizationName string) (*Organization, error)
	GetRepository(name string) (*Repository, error)
}

type GitCode struct {
	githubCode *GithubCode
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
