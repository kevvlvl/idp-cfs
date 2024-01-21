package git_api

import (
	"github.com/stretchr/testify/assert"
	"idp-cfs/global"
	"os"
	"testing"
)

func TestPushGoldenPath_ValidDirs_ValidGit_NoError(t *testing.T) {

	var (
		localReposFile = "./testdata/localTestRepos.tgz"
		tool           = global.ToolGithub
		codeUrl        = "./testdata/codeDir"
		codeUrlRemote  = "./testdata/codeRemoteDir"
		codeBranch     = "main"
		codeWorkdir    = "./testdata/idp-ut-code"
		gpUrl          = "./testdata/goldenPathDir"
		gpUrlRemote    = "./testdata/goldenPathRemoteDir"
		gpPathDir      = "."
		gpBranch       = "main"
		gpWorkDir      = "./testdata/idp-ut-gp"
	)

	r, err := os.Open(localReposFile)
	assert.Nil(t, err)

	global.ExtractTgz(r, "./testdata")

	err = pushGoldenPath(tool, codeUrl, codeBranch, gpUrl, gpPathDir, gpBranch, gpWorkDir, codeWorkdir, nil)

	// cleanup

	err = os.RemoveAll(codeUrl)
	assert.Nil(t, err)
	err = os.RemoveAll(codeUrlRemote)
	assert.Nil(t, err)
	err = os.RemoveAll(codeWorkdir)
	assert.Nil(t, err)
	err = os.RemoveAll(gpUrl)
	assert.Nil(t, err)
	err = os.RemoveAll(gpUrlRemote)
	assert.Nil(t, err)
	err = os.RemoveAll(gpWorkDir)
	assert.Nil(t, err)
}
