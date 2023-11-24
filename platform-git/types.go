package platform_git

import "github.com/google/go-github/v56/github"

type Git interface {
	Login()
	GetOrganizations() []*Organization
	GetRepositories() []*Repository
	GetBranches() []*Branch
}

type GitHub struct {
	Organizations []*github.Organization
	Repositories  []*github.Repository
}

type Organization struct {
	Name    string
	Company string
}

type Repository struct {
	Name  string
	Owner string
	URL   string
}

type Branch string
