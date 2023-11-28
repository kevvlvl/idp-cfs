package platform_git

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"os"
)

func GetGithubCode() *GitCode {

	pat := os.Getenv("GITHUB_PAT")
	client := github.NewClient(nil).WithAuthToken(pat)

	user, resp, err := client.Users.Get(context.Background(), "")

	valid := validateApiResponse(resp, err, "Error trying to get User")
	if !valid {
		return nil
	}

	return &GitCode{
		GithubClient: client,
		GithubUser:   user,
		Organization: nil,
		OrgExists:    false,
		Repository:   nil,
		RepoExists:   false,
	}
}

func GetGithubUser(c *github.Client) *github.User {

	user, resp, err := c.Users.Get(context.Background(), "")

	valid := validateApiResponse(resp, err, "Error trying to get User")
	if !valid {
		return nil
	}

	return user
}

func (c *GitCode) GetOrganization(organizationName string) (*Organization, error) {

	org, resp, err := c.GithubClient.Organizations.Get(context.Background(), organizationName)

	valid := validateApiResponse(resp, err, "Error trying to get organization")
	if !valid {
		return nil, errors.New(fmt.Sprintf("Error trying to get organization. Error: %v", err))
	}

	log.Debug().Msgf("Organization found %+v", org)

	return &Organization{
		Name:    org.Name,
		Company: org.Company,
	}, nil
}

func (c *GitCode) GetRepository(name string) (*Repository, error) {

	repo, resp, err := c.GithubClient.Repositories.Get(context.Background(), *c.GithubUser.Login, name)

	valid := validateApiResponse(resp, err, "Error trying to get repository")
	if !valid {
		return nil, errors.New(fmt.Sprintf("Error trying to get repository. Error: %v", err))
	}

	log.Debug().Msgf("Repository found %+v", repo)

	repoOrg := ""
	if repo.Organization != nil {
		repoOrg = *repo.Organization.Name
	}

	return &Repository{
		Name:         repo.Name,
		Organization: &repoOrg,
		Owner:        repo.Owner.Name,
		URL:          repo.URL,
	}, nil
}
