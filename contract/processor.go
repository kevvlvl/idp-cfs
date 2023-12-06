package contract

import (
	"errors"
	"github.com/rs/zerolog/log"
	"idp-cfs/platform_git"
	"idp-cfs/platform_gp"
	"strings"
)

func GetProcessor(contractFile string) *Processor {

	c, err := Load(contractFile)

	if err != nil {
		log.Error().Msg("ERROR when loading contract")
		return nil
	}

	var tag string
	if c.GoldenPath.Tag != nil {
		tag = *c.GoldenPath.Tag
	}

	gp := platform_gp.GetGoldenPath(*c.GoldenPath.Url,
		*c.GoldenPath.Name,
		*c.GoldenPath.Branch,
		*c.GoldenPath.Path,
		tag)

	return &Processor{
		Contract:   c,
		GitCode:    platform_git.GetCode(c.Code.Tool),
		GoldenPath: &gp,
	}
}

// Execute allows you to run the idp either in dryRun mode, or in real life (dryRun = false)
func (p *Processor) Execute(dryRunMode bool) (IdpStatus, error) {

	dryRunSuccess := performDryRun(p)

	if dryRunSuccess && dryRunMode {
		return IdpStatusSuccess, nil
	} else if dryRunSuccess {

		log.Info().Msg("Dry Run was successful. ***** BEGIN Creating platform... *****")

		if p.Contract.Action == NewContract {

			//----------------------------------------------------------------------------------
			// Create repository
			//----------------------------------------------------------------------------------

			err := p.GitCode.CreateRepository(p.Contract.Code.Repo)
			return evaluateIdpStatus(err)

		} else if p.Contract.Action == UpdateContract {

			//----------------------------------------------------------------------------------
			// Update repository
			//----------------------------------------------------------------------------------

			log.Error().Msg("Not implemented yet")
		}

		//----------------------------------------------------------------------------------
		// Push Golden Path into new or updated Repo
		//----------------------------------------------------------------------------------

		if p.Contract.GoldenPath.Url != nil {
			err := p.GitCode.PushFiles(p.GitCode.Repository, platform_gp.GetCheckoutPath())
			return evaluateIdpStatus(err)
		}

		return IdpStatusSuccess, nil
	} else {
		return IdpStatusFailure, nil
	}
}

func performDryRun(p *Processor) bool {

	//----------------------------------------------------------------------------------
	// GithubCode repository validation
	//----------------------------------------------------------------------------------

	code := platform_git.GetCode(p.Contract.Code.Tool)

	if p.Contract.Code.Org != nil {

		log.Info().Msgf("Contract Org defined. Search for %v", p.Contract.Code.Org)

		orgFound, err := code.GetOrganization(*p.Contract.Code.Org)

		if orgFound == nil && err != nil {
			return false
		}

		log.Info().Msgf("Org found: %v", orgFound)
	} else {
		log.Info().Msg("No Contract Org defined.")
	}

	//----------------------------------------------------------------------------------
	// GithubCode repository validation
	//----------------------------------------------------------------------------------
	// Search the code repo's organization

	repo, err := code.GetRepository(p.Contract.Code.Repo)

	if p.Contract.Action == NewContract {

		// In the case of a new infra request, we don't want to find an existing Git repo

		if err != nil {
			if strings.HasPrefix(err.Error(), "HTTP404") {
				// For the action new-contract, we want a HTTP 404! Otherwise, a new repo cannot be created.
				log.Info().Msg("New desired repo does not exist.")
			} else {
				log.Error().Msgf(err.Error())
				return false
			}
		} else {

			if repo != nil {
				repoFoundMsg := "A repository was found and returned. Make sure to review the code repo name and desired contract action"

				log.Error().Msgf(repoFoundMsg)
				return false
			} else {

				return false
			}
		}

	} else if p.Contract.Action == UpdateContract {

		// In the case of an update infra request, we want to find the repo and branch name
		if err != nil {
			log.Error().Msgf(err.Error())
			return false
		} else {

			if repo != nil {

				// For the action update-contract, we want an HTTP 2xx! Otherwise, no update can be done
				log.Info().Msgf("Found existing repo %v", repo)
			} else {
				return false
			}
		}

	} else {
		log.Error().Msgf("unexpected to get here. This means the contract was validated yet it made it here? Action: %v", p.Contract.Action)
		return false
	}

	//----------------------------------------------------------------------------------
	// Golden path section validation
	//----------------------------------------------------------------------------------
	if p.Contract.GoldenPath.Url != nil {

		err := p.GoldenPath.CloneGp()
		if err != nil {
			return false
		}

		err = platform_gp.DeleteClonePathDir()
		if err != nil {
			return false
		}

		log.Info().Msg("Checked out branch successfully.")
	}

	// TODO next: Kubernetes: create namespace and ensure namespace managed by ArgoCD

	// Verify kubernetes deployment section
	// Can I connect to k8s and verify the operator status?
	// do I have RBAC to create a namespace?
	// If logs is true, does grafana loki exist?

	// else if action == update-contract
	// call dry-run-update-contract func

	return true
}

func unexpectedError() error {
	m := "unexpected! Verify contract inputs, possibly report bug to the project maintainers"

	log.Error().Msg(m)
	return errors.New(m)
}

func evaluateIdpStatus(e error) (IdpStatus, error) {

	if e != nil {
		return IdpStatusFailure, e
	} else {
		return IdpStatusSuccess, nil
	}
}