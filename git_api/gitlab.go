package git_api

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"
	"idp-cfs/git_client"
	"idp-cfs/global"
	"os"
)

func (g *GitlabApi) ValidateNewCode(repoName string) error {

	project, err := getProject(g.client, repoName)

	if project == nil && err != nil {
		log.Info().Msgf("ValidateNewCode() - Project %s does not exist.", repoName)
		return nil
	} else if project != nil && err == nil {
		repoFound := fmt.Sprint("found Project when we did not expect one. Review contract code repo name and action")
		log.Warn().Msgf("ValidateNewCode() - %s", repoFound)
		return errors.New(repoFound)
	} else {
		log.Error().Msgf("ValidateNewCode() - Unexpected error returned: %v", err)
		return err
	}
}

func (g *GitlabApi) ValidateUpdateCode(repoName string) error {
	return nil
}

func (g *GitlabApi) ValidateGoldenPath(url, branch, workDir string) error {

	auth := getAuth(ToolGitlab)
	var gitAuth *git_client.GitClientAuth = nil

	if auth.gpDefined {
		gitAuth = git_client.GetAuth(auth.gpUser, auth.gpToken)
	}

	git := git_client.GetGitClient()
	_, err := git.CloneRepository(workDir, url, branch, gitAuth)
	if err != nil {
		log.Error().Msgf("ValidateGoldenPath() - Failed to clone golden path repo: %v", err)
	}

	log.Info().Msgf("ValidateGoldenPath() - Cloned the repo")

	// Delete the cloned repo if in dry-run. Otherwise, keep it to push this in the new code git repo
	err = os.RemoveAll(workDir)
	if err != nil {
		log.Error().Msgf("ValidateGoldenPath() - failed to delete the clone path: %v", err)
		return err
	}

	return nil
}

func (g *GitlabApi) CreateRepo(repoName string) error {

	p, err := createProject(g.client, repoName)
	if err != nil {
		log.Error().Msgf("CreateRepo() - Failed to Create Gitlab Project: %v", err)
		return err
	}

	g.project = p

	return nil
}

func (g *GitlabApi) PushGoldenPath(url, pathDir, branch, gpWorkdir, codeWorkDir string, tag *string) error {
	return pushGoldenPath(ToolGitlab, g.project.WebURL, g.project.DefaultBranch, url, pathDir, branch, gpWorkdir, codeWorkDir, tag)
}

func GetGitlabCodeClient(url string) *GitlabApi {

	auth := getAuth(ToolGitlab)
	if auth == nil {
		log.Error().Msg("GetGitlabCodeClient() - Cannot return client without auth info")
		return nil
	}

	if auth.codeDefined {
		return &GitlabApi{
			client: getClient(url, auth.codeToken),
		}
	}

	return nil
}

func GetGitlabGpClient(url string) *GitlabApi {

	auth := getAuth("gitlab")
	if auth == nil {
		log.Error().Msg("GetGitlabGpClient() - Cannot return client without auth info")
		return nil
	}

	if auth.gpDefined {
		return &GitlabApi{
			client: getClient(url, auth.gpToken),
		}
	}

	return nil
}

func getClient(url, token string) *gitlab.Client {

	c, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		log.Error().Msgf("getClient() - Failed to acquire Gitlab basic auth client: %v", err)
		return nil
	}

	return c
}

func getProject(g *gitlab.Client, projectName string) (*gitlab.Project, error) {

	p, resp, err := g.Projects.GetProject(projectName, nil)
	err = global.ValidateApiResponse(resp.Response, err, "Error trying to get project")
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("getProject() - Found Project: %s", p.Name)

	return p, nil
}

func createProject(g *gitlab.Client, projectName string) (*gitlab.Project, error) {

	log.Info().Msgf("createProject() - START Gitlab createProject: %s", projectName)

	p, err := getProject(g, projectName)

	if err == nil {
		msg := "found a Gitlab project with the name. Cannot create a new project"
		log.Warn().Msgf("createProject() - %s", msg)
		return nil, errors.New(msg)
	} else {

		opts := &gitlab.CreateProjectOptions{
			Name:          &projectName,
			DefaultBranch: global.StringPtr("main"),
		}

		newProject, resp, err := g.Projects.CreateProject(opts)
		err = global.ValidateApiResponse(resp.Response, err, "Error trying to create project")
		if err != nil {
			return nil, err
		}

		p = newProject
	}

	log.Info().Msgf("createProject() - END Gitlab createProject")
	return p, nil
}
