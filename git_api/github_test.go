package git_api

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/stretchr/testify/assert"
	"idp-cfs/global"
	"net/http"
	"testing"
)

func getTestRepoName() *string {
	return global.StringPtr("testRepo")
}

func getStubRepositoryWithOrg() *github.Repository {
	return &github.Repository{
		Name:         getTestRepoName(),
		Organization: getStubOrg(),
		Owner:        getStubUser(),
		URL:          global.StringPtr("http://cfs_idp_local_unit_test.local"),
	}
}

func getStubOrg() *github.Organization {
	return &github.Organization{
		Name: global.StringPtr("testOrg"),
	}
}

func getStubUser() *github.User {
	return &github.User{
		Name:  global.StringPtr("Unit McUnitTester"),
		Login: global.StringPtr("unittest"),
	}
}

func getStubValidResponse(code int) *github.Response {
	return &github.Response{
		Response: &http.Response{
			StatusCode: code,
		},
	}
}

func TestGetRepository_ValidRepoWithOrg_NoError(t *testing.T) {

	mockRepository := getStubRepositoryWithOrg()
	mockResponse := getStubValidResponse(200)

	c := getGithubClientWithoutAuth()
	c.user = getStubUser()
	c.getRepoFunc = func(ctx context.Context, c *github.Client, owner, repo string) (*github.Repository, *github.Response, error) {
		return mockRepository, mockResponse, nil
	}

	repository, err := c.getRepository(*getTestRepoName())
	assert.Nil(t, err)
	assert.NotNil(t, repository)
}

func TestHasAuth_Nil_Error(t *testing.T) {

	r := hasAuthUser(nil)
	assert.False(t, r)
}

func TestCreateRepository_Valid_NoError(t *testing.T) {

	c := getGithubClientWithoutAuth()
	c.user = getStubUser()
	c.createRepoFunc = func(ctx context.Context, c *github.Client, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
		return getStubRepositoryWithOrg(), getStubValidResponse(200), nil
	}

	c.createFileFunc = func(ctx context.Context, c *github.Client, owner, repo, path string, opts *github.RepositoryContentFileOptions) error {
		return nil
	}

	repo, err := c.createRepository("unitTestRepository")
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}
