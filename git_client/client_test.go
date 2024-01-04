package git_client

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGitClient_CloneRepository_ValidPublicUrl_NoErrors(t *testing.T) {

	g := GetGitClient()

	var (
		tmpDir = "/tmp/idp-cfs-unittest"
		url    = "https://github.com/kevvlvl/idp-cfs.git"
		branch = "main"
	)

	repository, err := g.CloneRepository(tmpDir, url, branch, nil)

	assert.Nil(t, err)
	assert.NotNil(t, repository)

	err = os.RemoveAll(tmpDir)
	assert.Nil(t, err)
}

func TestGitClient_CloneRepository_InvalidPublicUrl_AuthError(t *testing.T) {

	g := GetGitClient()

	var (
		tmpDir = "/tmp/idp-cfs-unittest"
		url    = "https://github.com/kevvlvl/this_repo_does_not_exist.git"
		branch = "main"
	)

	repository, err := g.CloneRepository(tmpDir, url, branch, nil)

	assert.Nil(t, repository)
	assert.Contains(t, err.Error(), "authentication")
	assert.NotNil(t, err)

	err = os.RemoveAll(tmpDir)
	assert.Nil(t, err)
}

func TestGetAuth_ValidCreds_NoErrors(t *testing.T) {

	user := "testUser"
	token := "test123456"
	s := GetAuth(user, token)

	assert.NotNil(t, s)
	assert.Equal(t, user, s.User)
	assert.Equal(t, token, s.Token)
}

func TestGetAuth_MissingCreds_Errors(t *testing.T) {

	user := ""
	token := ""
	s := GetAuth(user, token)

	assert.Nil(t, s)
}
