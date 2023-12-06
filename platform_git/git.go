package platform_git

import (
	"errors"
	"github.com/rs/zerolog/log"
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
	return errors.New("not implemented yet")
}
