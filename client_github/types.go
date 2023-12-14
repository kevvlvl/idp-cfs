package client_github

import (
	"flag"
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

var (
	GithubPrivateRepository = flag.Bool("private", false, "Will created repo be private.")
	GithubDescription       = flag.String("description", "", "Created by idp-cfs")
	GithubAutoInit          = flag.Bool("auto-init", false, "Pass true to create an initial commit with empty README.")
)
