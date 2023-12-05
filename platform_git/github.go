package platform_git

import (
	"context"
	"github.com/google/go-github/v56/github"
	"os"
)

func GetGithubCode() *GithubCode {

	pat := os.Getenv("CFS_GITHUB_PAT")
	client := github.NewClient(nil).WithAuthToken(pat)

	user, resp, err := client.Users.Get(context.Background(), "")

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

	org, resp, err := c.GithubClient.Organizations.Get(context.Background(), organizationName)

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

	repo, resp, err := c.GithubClient.Repositories.Get(context.Background(), *c.githubUser.Login, name)

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
