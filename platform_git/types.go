package platform_git

import "github.com/google/go-github/v56/github"

type GitHub struct {
	Organizations []*github.Organization
	Repositories  []*github.Repository
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

type Branch string
type GithubClient github.Client
type GithubUser github.User
