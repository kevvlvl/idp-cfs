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

	return &Repository{
		Name:         repo.Name,
		Organization: getOrganizationName(repo),
		Owner:        repo.Owner.Name,
		URL:          repo.CloneURL,
	}, nil
}

func (c *GithubCode) CreateRepository(repoName string, branch string) (*Repository, error) {

	newRepo := &github.Repository{
		Name:        &repoName,
		Private:     GithubPrivateRepository,
		Description: GithubDescription,
		AutoInit:    GithubAutoInit,
	}

	repo, resp, err := c.GithubClient.Repositories.Create(ctx, "", newRepo)

	err = validateApiResponse(resp, err, "Error trying to create repository")
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Created the repo successfully! Created on (timestamp): %v", repo.CreatedAt)

	emptyCommit := &github.RepositoryContentFileOptions{
		Message: github.String("Initial commit"),
		Content: []byte(""),
	}

	login := *c.githubUser.Login

	_, _, err = c.GithubClient.Repositories.CreateFile(ctx, login, repoName, "README.md", emptyCommit)
	if err != nil {
		log.Error().Msgf("Error creating a file for the empty commit. Error: %v", err)
		return nil, err
	}

	return &Repository{
		Name:         repo.Name,
		Organization: getOrganizationName(repo),
		Owner:        repo.Owner.Name,
		URL:          repo.CloneURL,
	}, nil
}

func getOrganizationName(r *github.Repository) *string {

	if r.Organization != nil {
		return r.Organization.Name
	} else {
		return nil
	}
}
