package platform_gp

import "github.com/go-git/go-git/v5"

type GoldenPath struct {
	Tool           string
	Name           string
	URL            string
	Branch         string
	Path           string
	Tag            string
	repository     *git.Repository
	GpCheckoutPath string
}

const (
	// GpGithub for the gp repository of type Gitlab
	GpGithub = "github"
	// GpGitlab for the gp repository of type Gitlab
	GpGitlab = "gitlab"
	// GpGitea for the gp repository of type Gitea
	GpGitea = "gitea"
	// DefaultCheckoutPath for the default git clone path for the GoldenPath if none set in the env var
	DefaultCheckoutPath = "/tmp/idp-cfs-gp"
)
