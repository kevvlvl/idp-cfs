package contract

import (
	"idp-cfs/client_git"
	"idp-cfs/client_github"
)

// FileReader is a wrapper interface to make os.ReadFile testable
type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type CfsFileReader struct{}

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
	CodeClonePath      string
	CodeGitBasicAuth   *client_git.GitBasicAuth
	GpClonePath        string
	Contract           *Contract
	GitClient          *client_git.GitClient
	GithubBasicAuth    *client_github.GithubBasicAuth
	GithubClient       *client_github.GithubClient
	GithubOrganization *client_github.Organization
	GithubRepository   *client_github.Repository
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
	// CodeClonePath is the path where we git clone the code path to prepare pushing the golden path (when defined)
	CodeClonePath = "/tmp/idp-cfs-code"
	// GoldenPathClonePath is the path where we git clone the code path to prepare pushing the golden path (when defined)
	GoldenPathClonePath = "/tmp/idp-cfs-gp"
)
