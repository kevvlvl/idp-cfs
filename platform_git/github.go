package platform_git

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"os"
)

func Login() *GithubClient {

	githubPat := os.Getenv("GITHUB_PAT")
	client := getClient(githubPat)

	return (*GithubClient)(client)
}

func (g *GithubClient) GetOrganization(name string) Organization {

	ctx := context.Background()

	org, resp, err := g.Organizations.Get(ctx, name)

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

func (g *GithubClient) GetRepositories(name string) Repository {

	ctx := context.Background()
	u := g.getUser()

	repo, resp, err := g.Repositories.Get(ctx, *u.Login, name)

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

func (g *GithubClient) getUser() *GithubUser {

	ctx := context.Background()

	user, resp, err := g.Users.Get(ctx, "")
	if err != nil {
		log.Error().Msgf("Error trying to get the User. Response = %v, Error = %v", resp, err)
		return nil
	}

	return (*GithubUser)(user)
}

func getClient(pat string) *github.Client {
	return github.NewClient(nil).WithAuthToken(pat)
}
