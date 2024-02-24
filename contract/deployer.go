package contract

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"idp-cfs/git_api"
	"idp-cfs/global"
	"os"
)

func GetState(dryRun bool, contractFilePath string) *State {

	c, err := Load(contractFilePath)
	if err != nil {
		return nil
	}

	state := &State{
		DryRun:   dryRun,
		Contract: c,
	}

	switch c.Code.Tool {
	case global.ToolGitlab:
		state.Code = git_api.GetGitlabCodeClient(*c.Code.Url)
	case global.ToolGithub:
		state.Code = git_api.GetGithubCodeClient()

	default:
		unexpectedResult(fmt.Sprintf("GetState() - Tool = %s", c.Code.Tool))
	}

	switch c.GoldenPath.Tool {
	case global.ToolGitlab:
		state.GoldenPath = git_api.GetGitlabGpClient(c.GoldenPath.Url)
	case global.ToolGithub:
		state.GoldenPath = git_api.GetGithubGpClient()
	default:
		unexpectedResult(fmt.Sprintf("GetState() - Tool = %s", c.GoldenPath.Tool))
	}

	return state
}

func (s *State) Deploy() (IdpStatus, error) {

	if err := validateState(s); err != nil {
		return IdpStatusFailure, err
	}

	if err := validateLocalStorageDirs(s); err != nil {
		return IdpStatusFailure, err
	}

	// DRY-RUN SECTION STARTS HERE
	log.Info().Msg("Deploy() - START Dry-Run")

	switch s.Contract.Action {
	case global.NewCode:
		if err := s.Code.ValidateNewCode(s.Contract.Code.Repo); err != nil {
			return IdpStatusFailure, err
		}
	case global.UpdateCode:
		if err := s.Code.ValidateUpdateCode(s.Contract.Code.Repo); err != nil {
			return IdpStatusFailure, err
		}
	default:
		unexpectedResult(fmt.Sprintf("Deploy() - Action = %s", s.Contract.Action))
	}

	if s.GoldenPath != nil {

		if err := s.GoldenPath.ValidateGoldenPath(s.Contract.GoldenPath.Url, s.Contract.GoldenPath.Branch, *s.Contract.GoldenPath.Workdir); err != nil {
			return IdpStatusFailure, err
		}
	}

	log.Info().Msg("Deploy() - COMPLETED Dry-Run")

	// if dryRun is false, deploy!
	if s.DryRun == false {

		log.Info().Msg("Deploy() - START Real Deployment")

		switch s.Contract.Action {
		case global.NewCode:

			// 1. Create the code repo
			if err := s.Code.CreateRepo(s.Contract.Code.Repo); err != nil {
				return IdpStatusFailure, err
			}

			// 2. If the Golden Path repo is defined, push that code in the code repo
			if err := s.Code.PushGoldenPath(s.Contract.GoldenPath.Url,
				s.Contract.GoldenPath.Path,
				s.Contract.GoldenPath.Branch,
				*s.Contract.GoldenPath.Workdir,
				*s.Contract.Code.Workdir,
				s.Contract.GoldenPath.Tag); err != nil {
				return IdpStatusFailure, err
			}

		case global.UpdateCode:

			// 1. Update the code repo by creating and pushing a new branch
			if err := s.Code.UpdateRepo(s.Contract.Code.Repo, "feat/gp-upgrade"); err != nil {
				return IdpStatusFailure, err
			}

			// 2. Push Golden Path in that newly created branch

			// 3. Create a PR from branch to main branch

		default:
			unexpectedResult(fmt.Sprintf("Deploy() - Action = %s", s.Contract.Action))
		}

		log.Info().Msg("Deploy() - COMPLETED Real Deployment")
	}

	return IdpStatusSuccess, nil
}

func validateState(s *State) error {
	if s.Contract == nil {
		return errors.New("contract cannot be nil. Ensure the YAML parses properly")
	}

	if s.Code == nil {
		return errors.New("code cannot be nil. Ensure the git source is implemented properly")
	}

	if s.GoldenPath == nil {
		return errors.New("golden path cannot be nil. Ensure the git source is implemented properly")
	}

	return nil
}

func validateLocalStorageDirs(s *State) error {
	if _, err := os.Stat(*s.Contract.GoldenPath.Workdir); !os.IsNotExist(err) {
		msg := fmt.Sprintf("path %s exists. Please delete or change path to a non-existing directory!", *s.Contract.GoldenPath.Workdir)
		log.Error().Msgf("validateLocalStorageDirs() - %s", msg)
		return errors.New(msg)
	}

	if _, err := os.Stat(*s.Contract.Code.Workdir); !os.IsNotExist(err) {

		msg := fmt.Sprintf("path %s exists. Please delete or change path to a non-existing directory!", *s.Contract.Code.Workdir)
		log.Error().Msgf("validateLocalStorageDirs() - %s", msg)
		return errors.New(msg)
	}

	return nil
}

func unexpectedResult(details string) {
	log.Error().Msgf("Unexpected error?! %s", details)
}
