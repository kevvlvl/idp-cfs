package contract

import (
	"idp-cfs/platform_git"
	"idp-cfs/platform_gp"
)

type Contract struct {
	Action string `yaml:"action"`

	Code struct {
		Tool   string  `yaml:"tool"`
		Org    *string `yaml:"org,omitempty"`
		Repo   string  `yaml:"repo"`
		Branch string  `yaml:"branch"`
	} `yaml:"code"`

	GoldenPath struct {
		Url    *string `yaml:"url,omitempty"`
		Name   *string `yaml:"name,omitempty"`
		Path   *string `yaml:"path,omitempty"`
		Branch *string `yaml:"branch,omitempty"`
		Tag    *string `yaml:"tag,omitempty"`
	} `yaml:"golden-path"`

	Deployment struct {
		Kubernetes struct {
			ClusterUrl string `yaml:"cluster-url"`
			Namespace  string `yaml:"namespace"`
			Logs       bool   `yaml:"logs"`
		} `yaml:"kubernetes"`
	} `yaml:"deployment"`
}

type Processor struct {
	Contract   *Contract
	GitCode    *platform_git.GitCode
	GoldenPath *platform_gp.GoldenPath
}

type IdpStatus int64

const (
	IdpStatusSuccess IdpStatus = iota
	IdpStatusFailure
	IdpStatusPartial
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
