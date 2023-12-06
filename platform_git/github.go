package platform_git

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"os"
)

var ctx = context.Background()

func GetGithubCode() *GithubCode {

	pat := os.Getenv("CFS_GITHUB_PAT")
	client := github.NewClient(nil).WithAuthToken(pat)

	user, resp, err := client.Users.Get(ctx, "")

	err = validateApiResponse(resp, err, "Error trying to get User")
	if err != nil {
		return nil
	}

	return &GithubCode{
		GithubClient: client,
		githubUser:   user,
	}
}

func (c *GithubCode) GetOrganization(organizationName string) (*Organization, error) {

	org, resp, err := c.GithubClient.Organizations.Get(ctx, organizationName)

	err = validateApiResponse(resp, err, "Error trying to get organization")
	if err != nil {
		return nil, err
	}

	return &Organization{
		Name:    org.Name,
		Company: org.Company,
	}, nil
}

func (c *GithubCode) GetRepository(name string) (*Repository, error) {

	repo, resp, err := c.GithubClient.Repositories.Get(ctx, *c.githubUser.Login, name)

	err = validateApiResponse(resp, err, "Error trying to get repository")
	if err != nil {
		return nil, err
	}

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

func (c *GithubCode) CreateRepository(name string) error {

	newRepo := &github.Repository{
		Name:        &name,
		Private:     GithubPrivateRepository,
		Description: GithubDescription,
		AutoInit:    GithubAutoInit,
	}

	repo, resp, err := c.GithubClient.Repositories.Create(ctx, "", newRepo)

	err = validateApiResponse(resp, err, "Error trying to create repository")
	if err != nil {
		return err
	}

	log.Info().Msgf("Created the repo successfully! Created on (timestamp): %v", repo.CreatedAt)

	return nil
}
