package rules

import (
	"idp-cfs/contract"
	"idp-cfs/platform_git"
	"idp-cfs/platform_gp"
)

func GetProcessor(c *contract.Contract, g *platform_git.GitCode, gp *platform_gp.GoldenPath) *Processor {

	return &Processor{
		Contract:   c,
		GitCode:    g,
		GoldenPath: gp,
	}
}

// DryRun returns true if the simulation run is successful.
// Verifies that all systems are up and return expected status codes
func (p *Processor) DryRun() bool {

	return false
}

func (p *Processor) Execute() (RuleResult, error) {

	return RuleResult(Failure), nil
}
