package platform_git

import "github.com/google/go-github/v56/github"

type GitCode struct {
	GithubClient *github.Client
	GithubUser   *github.User
	Organization *Organization
	OrgExists    bool
	Repository   *Repository
	RepoExists   bool
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
