package client_github

import (
	"github.com/google/go-github/v56/github"
)

type GithubClient struct {
	client *github.Client
	user   *github.User
}

type GithubBasicAuth struct {
	user  string
	token string
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
