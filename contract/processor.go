package contract

import (
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

	//----------------------------------------------------------------------------------
	// GithubCode repository validation
	//----------------------------------------------------------------------------------

	code := platform_git.GetCode(p.Contract.Code.Tool)

	if p.Contract.Code.Org != nil {

		log.Info().Msgf("Contract Org defined. Search for %v", p.Contract.Code.Org)

		orgFound, err := code.GetOrganization(*p.Contract.Code.Org)

		if orgFound == nil && err != nil {
			return IdpStatusFailure, err
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

				log.Info().Msg("new desired repo does not exist.")

				// For the action new-contract, we want a HTTP 404! Otherwise, a new repo cannot be created.
				if !dryRunMode {

					log.Info().Msg("create the repo...")

					newCodeRepo, err := p.GitCode.CreateRepository(p.Contract.Code.Repo, p.Contract.Code.Branch)

					if err != nil {
						return IdpStatusFailure, err
					}

					p.GitCode.Repository = newCodeRepo
				}

			} else {
				log.Error().Msgf(err.Error())

				if err != nil {
					return IdpStatusFailure, err
				}
			}
		} else {

			if repo != nil {
				repoFoundMsg := "repository was found and returned. Make sure to review the code repo name and desired contract action"

				log.Error().Msgf(repoFoundMsg)
				if err != nil {
					return IdpStatusFailure, err
				}
			} else {

				if err != nil {
					return IdpStatusFailure, err
				}
			}
		}

	} else if p.Contract.Action == UpdateContract {

		// In the case of an update infra request, we want to find the repo and branch name
		if err != nil {
			log.Error().Msgf(err.Error())
			if err != nil {
				return IdpStatusFailure, err
			}

		} else {

			if repo != nil {

				// For the action update-contract, we want an HTTP 2xx! Otherwise, no update can be done
				log.Info().Msgf("found existing repo %v", repo)
			} else {
				if err != nil {
					return IdpStatusFailure, err
				}
			}
		}

	} else {
		log.Error().Msgf("unexpected to get here. This means the contract was validated yet it made it here? Action: %v", p.Contract.Action)
		if err != nil {
			return IdpStatusFailure, err
		}
	}

	//----------------------------------------------------------------------------------
	// Golden path section validation
	//----------------------------------------------------------------------------------
	if p.Contract.GoldenPath.Url != nil {

		err := p.GoldenPath.CloneGp()
		if err != nil {
			return IdpStatusFailure, err
		}

		if dryRunMode {
			// Delete the cloned repo if in dry-run. Otherwise, keep it to push this in the new code git repo

			err := platform_gp.DeleteClonePathDir()
			if err != nil {
				return IdpStatusFailure, err
			}
		} else {
			//----------------------------------------------------------------------------------
			// Push Golden Path into new or updated Repo
			//----------------------------------------------------------------------------------

			err := p.GitCode.PushFiles(
				*p.GitCode.Repository.URL,
				p.Contract.Code.Branch,
				platform_gp.GetCheckoutPath(),
				p.GoldenPath.Path)

			if err != nil {
				return IdpStatusFailure, err
			}
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

	return IdpStatusSuccess, nil
}
