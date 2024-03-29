package git_api

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"idp-cfs/git_client"
	"idp-cfs/global"
	"os"
	"strings"
)

func (g *GithubApi) ValidateNewCode(repoName string) error {

	repo, err := g.getRepository(repoName)
	if repo == nil && err != nil && strings.HasPrefix(err.Error(), "HTTP404") {
		log.Info().Msgf("Repo %s does not exist.", repoName)
	} else if repo != nil {

		repoFound := fmt.Sprint("found Repository with same name. Review contract code repo name")
		log.Warn().Msgf(repoFound)
		return errors.New(repoFound)
	} else {

		log.Error().Msgf("Unexpected error returned: %v", err)
		return err
	}

	return nil
}

func (g *GithubApi) ValidateUpdateCode(repoName string) error {
	return nil
}

func (g *GithubApi) ValidateGoldenPath(url, branch, workDir string) error {

	auth := getAuth(ToolGithub)
	var gitAuth *git_client.GitClientAuth = nil

	if auth.gpDefined {
		gitAuth = git_client.GetAuth(auth.gpUser, auth.gpToken)
	}

	git := git_client.GetGitClient()
	_, err := git.CloneRepository(workDir, url, branch, gitAuth)
	if err != nil {
		log.Error().Msgf("Failed to clone golden path repo: %v", err)
	}

	log.Info().Msgf("Cloned the repo")

	// Delete the cloned repo if in dry-run. Otherwise, keep it to push this in the new code git repo
	err = os.RemoveAll(workDir)
	if err != nil {
		log.Error().Msgf("failed to delete the clone path: %v", err)
		return err
	}

	return nil
}

func (g *GithubApi) CreateRepo(repoName string) error {

	newCodeRepo, err := g.createRepository(repoName)
	g.repository = newCodeRepo
	if err != nil {
		return err
	}

	return nil
}

func (g *GithubApi) PushGoldenPath(url, pathDir, branch, gpWorkdir, codeWorkDir string, tag *string) error {
	return pushGoldenPath(ToolGithub, *g.repository.CloneURL, *g.repository.DefaultBranch, url, pathDir, branch, gpWorkdir, codeWorkDir, tag)
}

func GetGithubCodeClient(url string) *GithubApi {

	auth := getAuth(ToolGithub)

	if auth == nil && url != "" {
		log.Error().Msg("Cannot return client without auth info for on-premise Github")
		return nil
	}

	if auth == nil {
		return getGithubClientWithoutAuth()
	} else if auth.codeDefined {
		return getGithubClient(auth.codeToken)
	}

	return nil
}

func GetGithubGpClient(url string) *GithubApi {

	auth := getAuth(ToolGithub)

	if auth == nil && url != "" {
		log.Error().Msg("Cannot return client without auth info for on-premise Github")
		return nil
	}

	if auth == nil {
		return getGithubClientWithoutAuth()
	} else if auth.gpDefined {
		return getGithubClient(auth.gpToken)
	}

	return nil
}

func getGithubClient(authToken string) *GithubApi {

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(authToken)

	user, resp, err := client.Users.Get(ctx, "")
	err = global.ValidateApiResponse(resp.Response, err, "Error trying to get User")
	if err != nil {
		return nil
	}

	log.Debug().Msgf("Found User %v: ", user)

	return &GithubApi{
		client: client,
		user:   user,
		ctx:    ctx,
		getRepoFunc: func(ctx context.Context, c *github.Client, owner, repo string) (*github.Repository, *github.Response, error) {
			return c.Repositories.Get(ctx, owner, repo)
		},
	}
}

func getGithubClientWithoutAuth() *GithubApi {
	return &GithubApi{
		client: github.NewClient(nil),
		ctx:    context.Background(),
		getRepoFunc: func(ctx context.Context, c *github.Client, owner, repo string) (*github.Repository, *github.Response, error) {
			return c.Repositories.Get(ctx, owner, repo)
		},
	}
}

func (g *GithubApi) createRepository(repo string) (*github.Repository, error) {
	log.Info().Msgf("CreateRepository - Create the repo %s", repo)

	if !hasAuthUser(g.user) {
		return nil, errors.New("not authenticated")
	}

	newRepo := &github.Repository{
		Name:        &repo,
		Private:     global.BoolPtr(false),
		Description: global.StringPtr(""),
		AutoInit:    global.BoolPtr(false),
	}

	r, resp, err := g.client.Repositories.Create(g.ctx, "", newRepo)

	err = global.ValidateApiResponse(resp.Response, err, "Error trying to create repository")
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Created the repo. Created on (timestamp): %v", r.CreatedAt)

	emptyCommit := &github.RepositoryContentFileOptions{
		Message: github.String("Initial commit"),
		Content: []byte(""),
	}

	login := *g.user.Login

	_, _, err = g.client.Repositories.CreateFile(g.ctx, login, repo, "README.md", emptyCommit)
	if err != nil {
		log.Error().Msgf("Error creating a file for the empty commit: %v", err)
		return nil, err
	}

	return r, nil
}

func (g *GithubApi) getRepository(repoName string) (*github.Repository, error) {

	log.Info().Msgf("GetRepository - Search for %s", repoName)

	if !hasAuthUser(g.user) {
		return nil, errors.New("not authenticated")
	}

	repo, resp, err := g.getRepoFunc(g.ctx, g.client, *g.user.Login, repoName)
	err = global.ValidateApiResponse(resp.Response, err, "Error trying to get repository")
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Repo found")

	return repo, nil
}

func hasAuthUser(u *github.User) bool {

	if u == nil {
		log.Error().Msg("github user is null. Need authentication to call Github API")
		return false
	}

	return true
}
