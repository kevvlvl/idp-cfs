package git_client

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

func TestGitClient_PushFilesSrcDstExist_NoErrors(t *testing.T) {

	g := GetGitClient()

	// Mock git functions to return expected results = nil errors
	g.headFunc = getValidHeadFunc()
	g.workTreeFunc = getValidWorkTree()
	g.addGlobFunc = getValidGlobFunc()
	g.statusFunc = getValidStatusFunc()
	g.commitFunc = getValidCommitFunc()
	g.pushFunc = getValidPushFunc()

	r := git.Repository{}
	auth := GetAuth("testUser", "testToken")

	err := g.PushFiles(&r, auth)
	assert.Nil(t, err)
}

func TestGitClient_PushFilesSrcDstExistNoGitHead_Error(t *testing.T) {

	g := GetGitClient()

	// Mock git functions to return expected results = nil errors
	g.headFunc = func(r *git.Repository) (*plumbing.Reference, error) {
		return nil, errors.New("test error")
	}

	g.workTreeFunc = getValidWorkTree()
	g.addGlobFunc = getValidGlobFunc()
	g.statusFunc = getValidStatusFunc()
	g.commitFunc = getValidCommitFunc()
	g.pushFunc = getValidPushFunc()

	r := git.Repository{}
	auth := GetAuth("testUser", "testToken")

	err := g.PushFiles(&r, auth)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to return HEAD")
}

func TestGitClient_PushFilesSrcDstExistNoWorktree_Error(t *testing.T) {

	g := GetGitClient()

	// Mock git functions to return expected results = nil errors
	g.headFunc = getValidHeadFunc()

	g.workTreeFunc = func(r *git.Repository) (*git.Worktree, error) {
		return nil, errors.New("test error")
	}

	g.addGlobFunc = getValidGlobFunc()
	g.statusFunc = getValidStatusFunc()
	g.commitFunc = getValidCommitFunc()
	g.pushFunc = getValidPushFunc()

	r := git.Repository{}
	auth := GetAuth("testUser", "testToken")

	err := g.PushFiles(&r, auth)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to return worktree")
}

func TestGitClient_PushFilesSrcDstExistNoGlob_Error(t *testing.T) {

	g := GetGitClient()

	// Mock git functions to return expected results = nil errors
	g.headFunc = getValidHeadFunc()
	g.workTreeFunc = getValidWorkTree()

	g.addGlobFunc = func(w *git.Worktree, glob string) error {
		return errors.New("test error")
	}

	g.statusFunc = getValidStatusFunc()
	g.commitFunc = getValidCommitFunc()
	g.pushFunc = getValidPushFunc()

	r := git.Repository{}
	auth := GetAuth("testUser", "testToken")

	err := g.PushFiles(&r, auth)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to add . to git")
}

func TestGitClient_PushFilesSrcDstExistNoStatus_Error(t *testing.T) {

	g := GetGitClient()

	// Mock git functions to return expected results = nil errors
	g.headFunc = getValidHeadFunc()
	g.workTreeFunc = getValidWorkTree()
	g.addGlobFunc = getValidGlobFunc()

	g.statusFunc = func(w *git.Worktree) (git.Status, error) {
		return nil, errors.New("test error")
	}

	g.commitFunc = getValidCommitFunc()
	g.pushFunc = getValidPushFunc()

	r := git.Repository{}
	auth := GetAuth("testUser", "testToken")

	err := g.PushFiles(&r, auth)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to get status")
}

func TestGitClient_PushFilesSrcDstExistNoCommit_Error(t *testing.T) {

	g := GetGitClient()

	// Mock git functions to return expected results = nil errors
	g.headFunc = getValidHeadFunc()
	g.workTreeFunc = getValidWorkTree()
	g.addGlobFunc = getValidGlobFunc()
	g.statusFunc = getValidStatusFunc()

	g.commitFunc = func(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error) {
		return plumbing.Hash{}, errors.New("test error")
	}

	g.pushFunc = getValidPushFunc()

	r := git.Repository{}
	auth := GetAuth("testUser", "testToken")

	err := g.PushFiles(&r, auth)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to commit")
}

func TestGitClient_PushFilesSrcDstExistNoPush_Error(t *testing.T) {

	g := GetGitClient()

	// Mock git functions to return expected results = nil errors
	g.headFunc = getValidHeadFunc()
	g.workTreeFunc = getValidWorkTree()
	g.addGlobFunc = getValidGlobFunc()
	g.statusFunc = getValidStatusFunc()
	g.commitFunc = getValidCommitFunc()
	g.pushFunc = func(r *git.Repository, o *git.PushOptions) error {
		return errors.New("test error")
	}

	r := git.Repository{}
	auth := GetAuth("testUser", "testToken")

	err := g.PushFiles(&r, auth)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed for push commit")
}

func getValidHeadFunc() func(r *git.Repository) (*plumbing.Reference, error) {
	return func(r *git.Repository) (*plumbing.Reference, error) {
		p := plumbing.NewReferenceFromStrings("test", "toto")
		return p, nil
	}
}

func getValidWorkTree() func(r *git.Repository) (*git.Worktree, error) {
	return func(r *git.Repository) (*git.Worktree, error) {
		w := git.Worktree{}
		return &w, nil
	}
}

func getValidGlobFunc() func(w *git.Worktree, glob string) error {
	return func(w *git.Worktree, glob string) error {
		return nil
	}
}

func getValidStatusFunc() func(w *git.Worktree) (git.Status, error) {
	return func(w *git.Worktree) (git.Status, error) {
		s := git.Status{}
		return s, nil
	}
}

func getValidCommitFunc() func(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error) {
	return func(w *git.Worktree, msg string, opts *git.CommitOptions) (plumbing.Hash, error) {
		h := plumbing.NewHash("test")
		return h, nil
	}
}

func getValidPushFunc() func(r *git.Repository, o *git.PushOptions) error {
	return func(r *git.Repository, o *git.PushOptions) error {
		return nil
	}
}
