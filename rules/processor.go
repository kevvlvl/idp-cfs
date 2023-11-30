package rules

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
	"idp-cfs/platform_git"
	"strings"
)

func GetProcessor(contractFile string) *Processor {

	c, err := contract.Load(contractFile)

	if err != nil {
		log.Error().Msg("ERROR when loading contract")
		return nil
	}

	return &Processor{
		Contract:   c,
		GitCode:    platform_git.GetGithubCode(),
		GoldenPath: nil,
	}
}

// DryRun returns true if the simulation run is successful.
// Verifies that all systems are up and return expected status codes
func (p *Processor) DryRun() (bool, error) {

	//----------------------------------------------------------------------------------
	// Code repository validation
	//----------------------------------------------------------------------------------
	// Search the code repo's organization

	if p.Contract.Code.Org != nil {

		log.Info().Msgf("Contract Org defined. Search for %v", p.Contract.Code.Org)

		orgFound, err := p.GitCode.GetOrganization(*p.Contract.Code.Org)

		if orgFound == nil && err != nil {
			return false, err
		}

		log.Info().Msgf("Org found: %v", orgFound)
	} else {
		log.Info().Msg("No Contract Org defined.")
	}

	//----------------------------------------------------------------------------------
	// Code repository validation
	//----------------------------------------------------------------------------------
	// Search the code repo's organization

	repo, err := p.GitCode.GetRepository(p.Contract.Code.Repo)

	if p.Contract.Action == NewContract {

		// In the case of a new infra request, we don't want to find an existing Git repo

		if err != nil {
			if strings.HasPrefix(err.Error(), "HTTP404") {
				// For the action new-contract, we want a HTTP 404! Otherwise, a new repo cannot be created.
				log.Info().Msg("New desired repo does not exist.")
			} else {
				log.Error().Msgf(err.Error())
				return false, err
			}
		} else {

			if repo != nil {
				repoFoundMsg := "A repository was found and returned. Make sure to review the code repo name and desired contract action"

				log.Error().Msgf(repoFoundMsg)
				return false, errors.New(repoFoundMsg)
			} else {

				return false, unexpectedError()
			}
		}

	} else if p.Contract.Action == UpdateContract {

		// In the case of an update infra request, we want to find the repo and branch name
		if err != nil {
			log.Error().Msgf(err.Error())
			return false, err
		} else {

			if repo != nil {

				// For the action update-contract, we want an HTTP 2xx! Otherwise, no update can be done
				log.Info().Msgf("Found existing repo %v", repo)
			} else {
				return false, unexpectedError()
			}
		}

	} else {
		unexpectedMsg := fmt.Sprintf("unexpected to get here. This means the contract was validated yet it made it here? Action: %v", p.Contract.Action)

		log.Error().Msgf(unexpectedMsg)
		return false, errors.New(unexpectedMsg)
	}

	// Verify golden path section
	// Can I connect to the git repo of the gp?
	// Does the branch exist? If no, FAIL with reason. If yes, continue
	// Does the relative path exist. If no, FAIL with reason. If yes, continue
	// Does the name of the specified gp exist? If no, FAIL with reason. If yes, continue

	// Verify kubernetes deployment section
	// Can I connect to k8s and verify the operator status?
	// do I have RBAC to create a namespace?
	// If logs is true, does grafana loki exist?

	// else if action == update-contract
	// call dry-run-update-contract func

	return true, nil
}

func (p *Processor) Execute() (RuleResult, error) {

	return RuleResult(Failure), nil
}

func unexpectedError() error {
	m := "unexpected! null error and null repo struct. Review contract inputs and project owner for an update"

	log.Error().Msg(m)
	return errors.New(m)
}
