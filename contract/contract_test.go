package contract

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadValidContract_NoErrors(t *testing.T) {

	c, err := Load(getTestContractFilename())

	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestLoadInvalidContract_MissingRequiredField_Error(t *testing.T) {

	c, err := Load(getTestContractMissingToolFilename())

	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "contract metadata is not valid")
}

func TestValidateValidContract_ThenTrue(t *testing.T) {

	c := validContractNewPlatform()
	assert.True(t, validate(&c))
}

func TestValidateInvalidContactAction_ThenFalse(t *testing.T) {

	c := invalidContractActionDeletePlatform()
	assert.False(t, validate(&c))
}

func TestValidateInvalidContactCodeTool_ThenFalse(t *testing.T) {

	c := invalidContractCodeTool()
	assert.False(t, validate(&c))
}

func validContractNewPlatform() Contract {
	return Contract{
		Action: "new-code",
		Code: &Code{
			Tool:   "github",
			Repo:   "my-test-repo",
			Branch: "main",
		},
		GoldenPath: &GoldenPath{
			Url:    "http://github.local/some_test_url_gp",
			Path:   "./",
			Branch: "main",
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

func invalidContractActionDeletePlatform() Contract {

	c := validContractNewPlatform()
	c.Action = "delete-contract"
	return c
}

func invalidContractCodeTool() Contract {
	c := validContractNewPlatform()
	c.Code.Tool = "SuperGitServer"
	return c
}

func getTestContractFilename() string {
	return "./testdata/platform-order-test-new.yaml"
}

func getTestContractMissingToolFilename() string {
	return "./testdata/platform-order-test-new-missingtool.yaml"
}
