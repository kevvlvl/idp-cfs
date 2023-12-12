package contract

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateValidContract_ThenTrue(t *testing.T) {
	assert.True(t, validate(validContractNewPlatform()))
}

func TestValidateInvalidContactAction_ThenFalse(t *testing.T) {
	assert.False(t, validate(invalidContractActionDeletePlatform()))
}

func TestValidateInvalidContactCodeTool_ThenFalse(t *testing.T) {
	assert.False(t, validate(invalidContractCodeTool()))
}

func validContractNewPlatform() *Contract {
	return &Contract{
		Action: "new-contract",
		Code: struct {
			Tool   string  `yaml:"tool"`
			Org    *string `yaml:"org,omitempty"`
			Repo   string  `yaml:"repo"`
			Branch string  `yaml:"branch"`
		}{
			Tool:   "github",
			Org:    nil,
			Repo:   "my-test-repo",
			Branch: "main",
		},
		GoldenPath: struct {
			Url    *string `yaml:"url,omitempty"`
			Path   *string `yaml:"path,omitempty"`
			Branch *string `yaml:"branch,omitempty"`
			Tag    *string `yaml:"tag,omitempty"`
		}{
			Url:    getStrPointer("http://github.local/some_test_url_gp"),
			Path:   getStrPointer("./"),
			Branch: getStrPointer("main"),
			Tag:    nil,
		},
		Deployment: struct {
			Kubernetes struct {
				ClusterUrl string `yaml:"cluster-url"`
				Namespace  string `yaml:"namespace"`
				Logs       bool   `yaml:"logs"`
			} `yaml:"kubernetes"`
		}{
			Kubernetes: struct {
				ClusterUrl string `yaml:"cluster-url"`
				Namespace  string `yaml:"namespace"`
				Logs       bool   `yaml:"logs"`
			}{
				ClusterUrl: "k8s.cluster.unit.test.local",
				Namespace:  "apps",
				Logs:       true,
			},
		},
	}
}

func invalidContractActionDeletePlatform() *Contract {

	c := validContractNewPlatform()
	c.Action = "delete-contract"
	return c
}

func invalidContractCodeTool() *Contract {
	c := validContractNewPlatform()
	c.Code.Tool = "SuperGitServer"
	return c
}

func getStrPointer(s string) *string {
	return &s
}
