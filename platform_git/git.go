package platform_git

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func GetCode(tool string) *GitCode {

	var code *GitCode

	switch tool {

	case CodeGithub:
		code = &GitCode{githubCode: GetGithubCode()}
	case CodeGitea:
		log.Warn().Msg("Not implemented yet!")
	case CodeGitlab:
		log.Warn().Msg("Not implemented yet!")
	default:
		log.Error().Msgf("Unexpected code tool system which somehow passed validation! Tool: %v", tool)

	}

	return code
}

func (c *GitCode) GetOrganization(org string) (*Organization, error) {

	if c.githubCode != nil {
		return c.githubCode.GetOrganization(org)
	} else {
		return nil, errors.New("not implemented yet")
	}
}

func (c *GitCode) GetRepository(repo string) (*Repository, error) {

	if c.githubCode != nil {

		r, err := c.githubCode.GetRepository(repo)
		c.Repository = r

		return r, err
	} else {
		return nil, errors.New("not implemented yet")
	}
}

func (c *GitCode) CreateRepository(name string) error {
	if c.githubCode != nil {
		return c.githubCode.CreateRepository(name)
	} else {
		return errors.New("not implemented yet")
	}
}

func (c *GitCode) PushFiles(r *Repository, gpPath string) error {

	var files []string

	err := filepath.Walk(gpPath, func(filepath string, info os.FileInfo, err error) error {

		files = append(files, filepath)
		return nil
	})
	if err != nil {
		log.Error().Msgf("Failed to walk the directory %s. Error: %v", gpPath, err)
	}

	log.Info().Msgf("List of files to commit: %+v", files)

	// TODO use go-git to connect to the repository (r) and add the files (files) add, commit, and push

	return errors.New("not implemented yet")
}
