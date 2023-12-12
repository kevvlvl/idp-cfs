package contract

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"idp-cfs/platform_git"
	"idp-cfs/platform_gp"
	"strings"
)

func GetProcessor(contractFile string, gpCheckoutPath string, codeClonePath string) (*Processor, error) {

	c, err := Load(&ActualFileReader{}, contractFile)

	if err != nil {
		log.Error().Msgf("Error loading contract: %v", err)
		return nil, err
	}

	var tag string
	if c.GoldenPath.Tag != nil {
		tag = *c.GoldenPath.Tag
	}

	gp := platform_gp.GetGoldenPath(*c.GoldenPath.Url, *c.GoldenPath.Branch, *c.GoldenPath.Path, tag, gpCheckoutPath)

	return &Processor{
		Contract:   c,
		GitCode:    platform_git.GetCode(c.Code.Tool, codeClonePath),
		GoldenPath: &gp,
	}, nil
}

// Execute allows you to run the idp either in dryRun mode, or in real life (dryRun = false)
func (p *Processor) Execute(dryRunMode bool) (IdpStatus, error) {

	// Contract Code section
	if p.GitCode == nil {
		log.Error().Msg("Failed to obtain Git client for Code section.")
		return IdpStatusFailure, errors.New("did not obtain Git client for Code section")
	}

	err := validateContractCodeOrganization(p, p.GitCode)
	if err != nil {
		log.Error().Msg("Failed to validate contract code organization")
		return IdpStatusFailure, err
	}

	err = validateContractCodeRepo(dryRunMode, p, p.GitCode)
	if err != nil {
		log.Error().Msg("Failed to validate contract code repo")
		return IdpStatusFailure, err
	}

	// Golden Path Code section

	err = validateGoldenPath(dryRunMode, p)
	if err != nil {
		log.Error().Msg("Failed to validate the golden path")
		return IdpStatusFailure, err
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

func validateContractCodeOrganization(p *Processor, code *platform_git.GitCode) error {

	if p.Contract.Code.Org != nil {

		org, err := code.GetOrganization(*p.Contract.Code.Org)
		p.GitCode.Organization = org

		if org == nil && err != nil {
			return err
		}

	} else {
		log.Info().Msg("Contract does not define a Code organization.")
	}

	return nil
}

func validateContractCodeRepo(dryRunMode bool, p *Processor, code *platform_git.GitCode) error {

	repo, err := code.GetRepository(p.Contract.Code.Repo)
	p.GitCode.Repository = repo

	if p.Contract.Action == NewContract {

		// HTTP 404 when the request is to create a new repo is normal. we expect this response code
		if repo == nil && err != nil && strings.HasPrefix(err.Error(), "HTTP404") {

			log.Info().Msgf("Repo %s does not exist.", p.Contract.Code.Repo)

			if !dryRunMode {

				newCodeRepo, err := p.GitCode.CreateRepository(p.Contract.Code.Repo, p.Contract.Code.Branch)
				p.GitCode.Repository = newCodeRepo

				if err != nil {
					return err
				}
			}
		} else if repo != nil {

			repoFound := fmt.Sprint("found Repository when we did not expect one. Review contract code repo name")
			log.Warn().Msgf(repoFound)
			return errors.New(repoFound)
		} else {
			log.Error().Msgf("Unexpected error returned: %v", err)
			return err
		}
	} else if p.Contract.Action == UpdateContract {

		// expect HTTP 200 for the repo. It must exist
		if repo != nil && err == nil {

			log.Info().Msgf("Repo %s exists.", p.Contract.Code.Repo)

			if !dryRunMode {

				// TODO: create branch defined in the update contract. This is where we will push the gp in
			}
		}
	}

	return nil
}

func validateGoldenPath(dryRunMode bool, p *Processor) error {

	if p.Contract.GoldenPath.Url != nil {

		err := p.GoldenPath.CloneGp()
		if err != nil {
			log.Error().Msgf("failed to clone the repo: %v", err)
			return err
		}

		if dryRunMode {

			log.Info().Msg("Dry-Run mode enabled. Delete the golden path repo we just cloned.")
			// Delete the cloned repo if in dry-run. Otherwise, keep it to push this in the new code git repo

			err := p.GoldenPath.DeleteClonePathDir()
			if err != nil {
				log.Error().Msgf("failed to delete the clone path: %v", err)
				return err
			}
		} else {

			// Push Golden Path into new or updated Repo
			err := p.GitCode.PushFiles(*p.GitCode.Repository.URL, p.Contract.Code.Branch, p.GoldenPath.Path, p.GoldenPath.GpCheckoutPath)
			if err != nil {
				log.Error().Msgf("Failed to push files to the code repo: %v", err)
				return err
			}
		}

		log.Info().Msg("Checked out branch successfully.")
	}

	return nil
}
