package platform_git

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetUsername_EnvExists_NoErrors(t *testing.T) {

	user := "test-user-123"
	err := os.Setenv("CFS_GITHUB_USER", user)

	if err != nil {
		t.Fatalf("Failed to set env for unit test: %v", err)
	}

	returnedUser := GetUsername()

	assert.NotEmpty(t, returnedUser)
	assert.Equal(t, user, returnedUser, "The returned username from env var does not match os.Setenv")
}

func TestGetUsername_EnvEmpty_EmptyStr(t *testing.T) {

	user := ""
	err := os.Setenv("CFS_GITHUB_USER", user)

	if err != nil {
		t.Fatalf("Failed to set env for unit test: %v", err)
	}

	returnedUser := GetUsername()

	assert.Empty(t, returnedUser)
}

func TestGetPersonalAccessToken_EnvExists_NoErrors(t *testing.T) {

	pat := "abcdef11112233344455666777unittesttoken"
	err := os.Setenv("CFS_GITHUB_PAT", pat)

	if err != nil {
		t.Fatalf("Failed to set env for unit test: %v", err)
	}

	returnedPat := GetPersonalAccessToken()

	assert.NotEmpty(t, returnedPat)
	assert.Equal(t, pat, returnedPat, "The returned username from env var does not match os.Setenv")
}

func TestGetPersonalAccessToken_EnvEmpty_EmptyStr(t *testing.T) {

	pat := ""
	err := os.Setenv("CFS_GITHUB_PAT", pat)

	if err != nil {
		t.Fatalf("Failed to set env for unit test: %v", err)
	}

	returnedPat := GetPersonalAccessToken()

	assert.Empty(t, returnedPat)
}
