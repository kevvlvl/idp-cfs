package rules

import (
	"idp-cfs/contract"
	"idp-cfs/platform_git"
	"idp-cfs/platform_gp"
)

type Processor struct {
	Contract   *contract.Contract
	GitCode    *platform_git.GitCode
	GoldenPath *platform_gp.GoldenPath
}

type RuleResult int64

const (
	Success RuleResult = iota
	Failure
	Partial
)
