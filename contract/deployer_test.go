package contract

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"idp-cfs2/git_api"
	"idp-cfs2/global"
	"os"
	"testing"
)

func TestValidateState_NoError(t *testing.T) {

	s := getDryRunGithubState()
	err := validateState(s)
	assert.Nil(t, err)
}

func TestValidateState_NilContract_Error(t *testing.T) {

	s := getDryRunGithubState()
	s.Contract = nil
	err := validateState(s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "contract cannot be nil")
}

func TestValidateState_NilCode_Error(t *testing.T) {

	s := getDryRunGithubState()
	s.Code = nil
	err := validateState(s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "code cannot be nil")
}

func TestValidateState_NilGoldenPath_Error(t *testing.T) {

	s := getDryRunGithubState()
	s.GoldenPath = nil
	err := validateState(s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "golden path cannot be nil")
}

func TestValidateLocalStorageDirs_NoDirs_NoErrors(t *testing.T) {

	s := getDryRunGithubState()

	s.Contract.GoldenPath = &GoldenPath{
		Workdir: global.StringPtr("/tmp/test_gp_dir"),
	}

	s.Contract.Code = &Code{
		Workdir: global.StringPtr("/tmp/test_code_dir"),
	}

	// test
	err := validateLocalStorageDirs(s)
	assert.Nil(t, err)
}

func TestValidateLocalStorageDirs_CodeExists_Error(t *testing.T) {

	s := getDryRunGithubState()

	s.Contract.GoldenPath = &GoldenPath{
		Workdir: global.StringPtr("/tmp/test_gp_dir"),
	}

	s.Contract.Code = &Code{
		Workdir: global.StringPtr("/tmp/test_code_dir"),
	}

	// test data create
	err := global.CreateFolder(*s.Contract.Code.Workdir)
	assert.Nil(t, err)

	// test
	err = validateLocalStorageDirs(s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("path %s exists.", *s.Contract.Code.Workdir))

	// test data cleanup
	err = os.RemoveAll(*s.Contract.Code.Workdir)
	assert.Nil(t, err)
}

func TestValidateLocalStorageDirs_GoldenPathExists_Error(t *testing.T) {

	s := getDryRunGithubState()

	s.Contract.GoldenPath = &GoldenPath{
		Workdir: global.StringPtr("/tmp/test_gp_dir"),
	}

	s.Contract.Code = &Code{
		Workdir: global.StringPtr("/tmp/test_code_dir"),
	}

	// test data create
	err := global.CreateFolder(*s.Contract.GoldenPath.Workdir)
	assert.Nil(t, err)

	// test
	err = validateLocalStorageDirs(s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("path %s exists.", *s.Contract.GoldenPath.Workdir))

	// test data cleanup
	err = os.RemoveAll(*s.Contract.GoldenPath.Workdir)
	assert.Nil(t, err)
}

func getDryRunGithubState() *State {
	return &State{
		DryRun:     true,
		Contract:   &Contract{},
		Code:       &git_api.GithubApi{},
		GoldenPath: &git_api.GithubApi{},
	}
}
