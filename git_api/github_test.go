package git_api

import (
	"context"
	"github.com/google/go-github/v56/github"
	"github.com/stretchr/testify/assert"
	"idp-cfs2/util"
	"net/http"
	"testing"
)

func getTestRepoName() *string {
	return util.StringPtr("testRepo")
}

func getStubRepositoryWithOrg() *github.Repository {
	return &github.Repository{
		Name:         getTestRepoName(),
		Organization: getStubOrg(),
		Owner:        getStubUser(),
		URL:          util.StringPtr("http://cfs_idp_local_unit_test.local"),
	}
}

func getStubOrg() *github.Organization {
	return &github.Organization{
		Name: util.StringPtr("testOrg"),
	}
}

func getStubUser() *github.User {
	return &github.User{
		Name:  util.StringPtr("Unit McUnitTester"),
		Login: util.StringPtr("unittest"),
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
