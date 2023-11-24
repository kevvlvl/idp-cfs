package platform_git

import "github.com/google/go-github/v56/github"

type GithubClient struct {
	Client *github.Client
	User   *github.User
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
