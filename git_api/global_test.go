package git_api

import (
	"testing"
)

func TestPushGoldenPath_ValidDirs_ValidGit_NoError(t *testing.T) {

	//var (
	//	tool          = global.ToolGithub
	//	codeUrl       = "./testdata/codeDir"
	//	codeUrlRemote = "./testdata/codeRemoteDir"
	//	codeBranch    = "main"
	//	codeWorkdir   = "./testdata/idp-ut-code"
	//	gpUrl         = "./testdata/goldenPathDir"
	//	gpUrlRemote   = "./testdata/goldenPathRemoteDir"
	//	gpPathDir     = "."
	//	gpBranch      = "main"
	//	gpWorkDir     = "./testdata/idp-ut-gp"
	//)
	//
	// Revise best way to prepare git test data. TestContainers probably
	//// create testdata
	//err := global.CreateFolder(codeUrl)
	//assert.Nil(t, err)
	//
	//err = global.CreateFolder(codeUrlRemote)
	//assert.Nil(t, err)
	//
	//err = global.CreateFolder(gpUrl)
	//assert.Nil(t, err)
	//
	//err = global.CreateFolder(gpUrlRemote)
	//assert.Nil(t, err)
	//
	//_, err = os.Create(path.Join(gpUrl, "code.txt"))
	//assert.Nil(t, err)
	//
	//// cleanup test data
	//err = os.RemoveAll(codeUrl)
	//assert.Nil(t, err)
	//err = os.RemoveAll(codeUrlRemote)
	//assert.Nil(t, err)
	//err = os.RemoveAll(codeWorkdir)
	//assert.Nil(t, err)
	//err = os.RemoveAll(gpUrl)
	//assert.Nil(t, err)
	//err = os.RemoveAll(gpUrlRemote)
	//assert.Nil(t, err)
	//err = os.RemoveAll(gpWorkDir)
	//assert.Nil(t, err)
}
