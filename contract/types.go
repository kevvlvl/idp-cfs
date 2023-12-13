package contract

import (
	"idp-cfs/platform_git"
	"idp-cfs/platform_gp"
)

// FileReader is a wrapper interface to make os.ReadFile testable
type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type ActualFileReader struct{}

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
)
