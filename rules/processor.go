package rules

import (
	"github.com/rs/zerolog/log"
	"idp-cfs/contract"
	"idp-cfs/platform_git"
)

func GetProcessor(c *contract.Contract) *Processor {

	git := platform_git.GetGithubCode()

	return &Processor{
		Contract:   c,
		GitCode:    git,
		GoldenPath: nil,
	}
}

// DryRun returns true if the simulation run is successful.
// Verifies that all systems are up and return expected status codes
func (p *Processor) DryRun() bool {

	switch p.Contract.Action {
	case "new-contract":

		orgFound, err := p.GitCode.GetOrganization(p.Contract.Code.Org)

		if err != nil {
			log.Warn().Msgf(err.Error())
		}

		log.Info().Msgf("Org found: %v", orgFound)

	case "update-contract":
	default:
		log.Error().Msgf("Unexpected to get here. This means the contract was validated yet it made it here? Action: %v", p.Contract.Action)
	}

	if p.Contract.Action == "new-contract" {

	}

	// if action == new-contract

	// verify Code section:
	// Can I connect to git?
	// Can I find the org (if any defined)
	// Can I find a repo with the same name? If yes. FAIL with reason. If not, continue

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

	//code.Organization, code.OrgExists = code.GetOrganization("kevvlvl")
	//code.Repository, code.RepoExists = code.GetRepository("idp-cfs")

	return false
}

func (p *Processor) Execute() (RuleResult, error) {

	return RuleResult(Failure), nil
}
