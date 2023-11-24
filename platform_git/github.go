package platform_git

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"os"
)

func Github() GithubClient {

	githubPat := os.Getenv("GITHUB_PAT")
	client := getClient(githubPat)

	return GithubClient{
		Client: client,
		User:   getUser(client),
	}
}

func (g *GithubClient) GetOrganization(name string) Organization {

	ctx := context.Background()

	org, resp, err := g.Client.Organizations.Get(ctx, name)

	valid := validateApiResponse(resp, err, "Error trying to get organization")
	if !valid {
		return Organization{}
	}

	log.Debug().Msgf("Organization found %+v", org)

	return Organization{
		Name:    org.Name,
		Company: org.Company,
	}
}

func (g *GithubClient) GetRepository(name string) Repository {

	ctx := context.Background()

	repo, resp, err := g.Client.Repositories.Get(ctx, *g.User.Login, name)

	valid := validateApiResponse(resp, err, "Error trying to get Repository")
	if !valid {
		return Repository{}
	}

	log.Debug().Msgf("Repository found %+v", repo)

	repoOrg := ""
	if repo.Organization != nil {
		repoOrg = *repo.Organization.Name
	}

	return Repository{
		Name:         repo.Name,
		Organization: &repoOrg,
		Owner:        repo.Owner.Name,
		URL:          repo.URL,
	}
}

func getUser(c *github.Client) *github.User {

	ctx := context.Background()

	user, resp, err := c.Users.Get(ctx, "")
	if err != nil {
		log.Error().Msgf("Error trying to get the User. Response = %v, Error = %v", resp, err)
		return nil
	}

	return user
}

func getClient(pat string) *github.Client {
	return github.NewClient(nil).WithAuthToken(pat)
}
