package contract

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"idp-cfs/client_git"
	"idp-cfs/client_github"
	"idp-cfs/util"
	"os"
	"path"
	"strings"
)

func GetProcessor(contractFile string, gpClonePath string, codeClonePath string) (*Processor, error) {

	c, err := Load(&CfsFileReader{}, contractFile)

	if err != nil {
		log.Error().Msgf("Error loading contract: %v", err)
		return nil, err
	}

	user := os.Getenv("CFS_CODE_GIT_USER")
	token := os.Getenv("CFS_CODE_GIT_PAT")

	if c.Code.Tool == client_git.CodeGithub {
		return &Processor{
			CodeClonePath:    codeClonePath,
			CodeGitBasicAuth: client_git.GetAuth(user, token),
			Contract:         c,
			GitClient:        client_git.GetGitClient(),
			GithubBasicAuth:  client_github.GetAuth(user, token),
			GithubClient:     client_github.GetGithubClient(client_github.GetAuth(user, token)),
			GpClonePath:      gpClonePath,
		}, nil
	} else {
		return nil, errors.New("code other than Github not implemented yet")
	}
}

// Execute allows you to run the idp either in dryRun mode, or in real life (dryRun = false)
func (p *Processor) Execute(dryRunMode bool) (IdpStatus, error) {

	// Contract Code section
	if p.GitClient == nil {
		log.Error().Msg("Failed to obtain Git client for Code section.")
		return IdpStatusFailure, errors.New("did not obtain Git client for Code section")
	}

	err := p.validateContractCodeOrganization()
	if err != nil {
		log.Error().Msg("Failed to validate contract code organization")
		return IdpStatusFailure, err
	}

	err = p.validateContractCodeRepo(dryRunMode)
	if err != nil {
		log.Error().Msg("Failed to validate contract code repo")
		return IdpStatusFailure, err
	}

	// Golden Path Code section

	err = p.validateGoldenPath(dryRunMode)
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

// validateContractCodeOrganization verifies and loads the git organization if any, and sets it to the Processor
func (p *Processor) validateContractCodeOrganization() error {

	if p.Contract.Code.Tool == client_git.CodeGithub && p.Contract.Code.Org != nil {
		org, err := p.GithubClient.GetOrganization(*p.Contract.Code.Org)

		if org == nil && err != nil {
			return err
		}

		p.GithubOrganization = org
	}

	return nil
}

func (p *Processor) validateContractCodeRepo(dryRunMode bool) error {

	if p.Contract.Code.Tool == client_git.CodeGithub {
		repo, err := p.GithubClient.GetRepository(p.Contract.Code.Repo)

		if p.Contract.Action == NewContract {

			// HTTP 404 when the request is to create a new repo is normal. we expect this response code
			if repo == nil && err != nil && strings.HasPrefix(err.Error(), "HTTP404") {

				log.Info().Msgf("Repo %s does not exist.", p.Contract.Code.Repo)

				if !dryRunMode {

					newCodeRepo, err := p.GithubClient.CreateRepository(p.Contract.Code.Repo)
					p.GithubRepository = newCodeRepo

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

				p.GithubRepository = repo
				log.Info().Msgf("Repo %s exists.", p.Contract.Code.Repo)

				if !dryRunMode {

				}
			}
		}
	}

	return nil
}

func (p *Processor) validateGoldenPath(dryRunMode bool) error {

	if p.Contract.GoldenPath.Url != nil {

		_, err := p.GitClient.CloneRepository(p.GpClonePath, *p.Contract.GoldenPath.Url, p.Contract.GoldenPath.Branch, nil)

		if err != nil {
			log.Error().Msgf("failed to clone the repo: %v", err)
			return err
		}

		if dryRunMode {

			log.Info().Msg("Dry-Run mode enabled. Delete the golden path repo we just cloned.")
			// Delete the cloned repo if in dry-run. Otherwise, keep it to push this in the new code git repo

			err := util.RemoveAllDir(p.GpClonePath)
			if err != nil {
				log.Error().Msgf("failed to delete the clone path: %v", err)
				return err
			}
		} else {

			// GP_S1. Clone the newly created code repo
			codeRepo, err := p.GitClient.CloneRepository(p.CodeClonePath, *p.GithubRepository.URL, &p.Contract.Code.Branch, p.CodeGitBasicAuth)
			if err != nil {
				log.Error().Msgf("Failed to clone the code repository: %v", err)
				return err
			}

			// GP_S2. Copy files from the cloned GP to the cloned code repo
			var gpPath string
			if p.Contract.GoldenPath.Path == nil {
				gpPath = p.GpClonePath
			} else {
				gpPath = path.Join(p.GpClonePath, *p.Contract.GoldenPath.Path)
			}

			err = util.CopyFilesDeep(gpPath, p.CodeClonePath)
			if err != nil {
				log.Error().Msgf("Failed to copy files from goldenpath to code repo: %v", err)
				return err
			}

			// GP_S3. Push GP files to code repo
			err = p.GitClient.PushFiles(codeRepo, p.CodeClonePath, p.CodeGitBasicAuth)
			if err != nil {
				log.Error().Msgf("Failed to push files to the code repo: %v", err)
				return err
			}

		}

		log.Info().Msg("Checked out branch successfully.")
	}

	return nil
}
