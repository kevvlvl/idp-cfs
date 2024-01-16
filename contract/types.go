package contract

import "idp-cfs2/git_api"

type Code struct {
	Tool    string  `yaml:"tool"`
	Url     *string `yaml:"url,omitempty"`
	Repo    string  `yaml:"repo"`
	Branch  string  `yaml:"branch"`
	Workdir *string `yaml:"workdir,omitempty"`
}

type GoldenPath struct {
	Tool    string  `yaml:"tool"`
	Url     string  `yaml:"url"`
	Path    string  `yaml:"path"`
	Branch  string  `yaml:"branch"`
	Tag     *string `yaml:"tag,omitempty"`
	Workdir *string `yaml:"workdir,omitempty"`
}

type Contract struct {
	Action     string      `yaml:"action"`
	Code       *Code       `yaml:"code"`
	GoldenPath *GoldenPath `yaml:"golden-path,omitempty"`

	Deployment struct {
		Kubernetes struct {
			ClusterUrl string `yaml:"cluster-url"`
			Namespace  string `yaml:"namespace"`
			Logs       bool   `yaml:"logs"`
		} `yaml:"kubernetes"`
	} `yaml:"deployment"`
}

type State struct {
	DryRun     bool
	Contract   *Contract
	Code       git_api.GitSource
	GoldenPath git_api.GitSource
}

type IdpStatus int64

const (
	IdpStatusSuccess IdpStatus = iota
	IdpStatusFailure
	IdpStatusPartial
)
