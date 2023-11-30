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

const (
	// NewContract for request of New infrastructure
	NewContract = "new-contract"
	// UpdateContract for request to Update an existing infrastructure
	UpdateContract = "update-contract"
	// CodeGithub for the code repository of type Github (public/cloud)
	CodeGithub = "github"
	// CodeGitlab for the code repository of type Gitlab
	CodeGitlab = "gitlab"
	// CodeGitea for the code repository of type Gitea
	CodeGitea = "gitea"
)
