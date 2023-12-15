package client_github

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"idp-cfs/util"
)

var ctx = context.Background()

func GetGithubClient(auth *GithubBasicAuth) *GithubClient {

	log.Info().Msg("GetGithubClient - Init a Github Client")

	var client *github.Client

	if auth == nil {
		client = github.NewClient(nil)
	} else {
		client = github.NewClient(nil).WithAuthToken(auth.token)
	}

	user, resp, err := client.Users.Get(ctx, "")

	err = validateApiResponse(resp, err, "Error trying to get User")
	if err != nil {
		return nil
	}

	return &GithubClient{
		client: client,
		user:   user,
	}
}

func (g *GithubClient) GetRepository(repoName string) (*Repository, error) {

	log.Info().Msgf("GetRepository - Search for %s", repoName)

	repo, resp, err := g.client.Repositories.Get(ctx, *g.user.Login, repoName)

	err = validateApiResponse(resp, err, "Error trying to get repository")
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Repo found")

	return &Repository{
		Name:         repo.Name,
		Organization: getOrganizationName(repo),
		Owner:        repo.Owner.Name,
		URL:          repo.CloneURL,
	}, nil
}

func (g *GithubClient) GetOrganization(organizationName string) (*Organization, error) {

	log.Info().Msgf("GetOrganization - Search for %s", organizationName)

	org, resp, err := g.client.Organizations.Get(ctx, organizationName)

	err = validateApiResponse(resp, err, "Error trying to get organization")
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Organization found")

	return &Organization{
		Name:    org.Name,
		Company: org.Company,
	}, nil
}

func (g *GithubClient) CreateRepository(repo string) (*Repository, error) {
	log.Info().Msgf("CreateRepository - Create the repo %s", repo)

	newRepo := &github.Repository{
		Name:        &repo,
		Private:     util.BoolPtr(false),
		Description: util.StringPtr(""),
		AutoInit:    util.BoolPtr(false),
	}

	r, resp, err := g.client.Repositories.Create(ctx, "", newRepo)

	err = validateApiResponse(resp, err, "Error trying to create repository")
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Created the repo. Created on (timestamp): %v", r.CreatedAt)

	emptyCommit := &github.RepositoryContentFileOptions{
		Message: github.String("Initial commit"),
		Content: []byte(""),
	}

	login := *g.user.Login

	_, _, err = g.client.Repositories.CreateFile(ctx, login, repo, "README.md", emptyCommit)
	if err != nil {
		log.Error().Msgf("Error creating a file for the empty commit: %v", err)
		return nil, err
	}

	return &Repository{
		Name:         r.Name,
		Organization: getOrganizationName(r),
		Owner:        r.Owner.Name,
		URL:          r.CloneURL,
	}, nil
}

func validateApiResponse(resp *github.Response, e error, msg string) error {

	if e != nil {
		if resp.Response.StatusCode == 404 {

			errorMsg := fmt.Sprintf("HTTP404 - "+msg+" - response: %v - error: %v", resp, e)
			log.Warn().Msgf(errorMsg)
			return errors.New(errorMsg)

		} else if resp.Response.StatusCode >= 400 && resp.Response.StatusCode <= 499 {

			errorMsg := fmt.Sprintf("HTTP4xx - "+msg+" - response: %v - error: %v", resp, e)

			log.Warn().Msgf(errorMsg)
			return errors.New(errorMsg)
		} else if resp.Response.StatusCode >= 500 && resp.Response.StatusCode <= 599 {

			errorMsg := fmt.Sprintf("HTTP5xx - "+msg+" - response: %v - error: %v", resp, e)

			log.Error().Msgf(errorMsg)
			return errors.New(errorMsg)
		}
	}

	return nil
}

func GetAuth(user string, token string) *GithubBasicAuth {

	if user == "" && token == "" {
		log.Debug().Msg("getAuth - No basic auth env defined.")
		return nil
	} else {
		log.Debug().Msg("getAuth - Basic auth env defined")
		return &GithubBasicAuth{
			user:  user,
			token: token,
		}
	}
}

func getOrganizationName(r *github.Repository) *string {

	if r.Organization != nil {
		return r.Organization.Name
	} else {
		return nil
	}
}
