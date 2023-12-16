package client_github

import (
	"errors"
	"github.com/google/go-github/v56/github"
	"github.com/stretchr/testify/assert"
	"idp-cfs/util"
	"net/http"
	"testing"
)

func TestValidateApiResponseWithNilError_ThenNil(t *testing.T) {

	resultErr := validateApiResponse(respClientError(), nil, "Client error")

	assert.Nil(t, resultErr)
}

func TestValidateApiResponseWithClientError_ThenErrorHttp4xx(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := validateApiResponse(respClientError(), err, msg)

	assert.NotNil(t, resultErr)
	assert.Contains(t, resultErr.Error(), "HTTP4xx - ")
}

func TestValidateApiResponseWithClientError_ThenErrorHttp404(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := validateApiResponse(respClientNotFoundError(), err, msg)

	assert.NotNil(t, resultErr)
	assert.Contains(t, resultErr.Error(), "HTTP404 - ")
}

func TestValidateApiResponseWithClientError_ThenErrorHttp5xx(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := validateApiResponse(respServerError(), err, msg)

	assert.NotNil(t, resultErr)
	assert.Contains(t, resultErr.Error(), "HTTP5xx - ")
}

func TestGetAuth_ValidCreds_NoErrors(t *testing.T) {

	user := "testUser"
	token := "test123456"
	s := GetAuth(user, token)

	assert.NotNil(t, s)
	assert.Equal(t, user, s.user)
	assert.Equal(t, token, s.token)
}

func TestGetAuth_MissingCreds_Error(t *testing.T) {
	user := ""
	token := ""
	s := GetAuth(user, token)

	assert.Nil(t, s)
}

func TestGetOrganizationName_NotNil_GetName(t *testing.T) {

	testOrgName := util.StringPtr("TestOrganization")

	r := github.Repository{
		Organization: &github.Organization{
			Name: testOrgName,
		},
	}

	orgName := getOrganizationName(&r)

	assert.NotNil(t, orgName)
	assert.Equal(t, testOrgName, r.Organization.Name)
}

func errorTestNotNil() (error, string) {

	errorMsg := "failed to query github"
	return errors.New(errorMsg), errorMsg
}

func respClientError() *github.Response {
	return getGithubResponse("Bad Request", 400)
}

func respClientNotFoundError() *github.Response {
	return getGithubResponse("Not Found", 404)
}

func respServerError() *github.Response {
	return getGithubResponse("Internal Server Error", 500)
}

func getGithubResponse(status string, code int) *github.Response {
	return &github.Response{
		Response: &http.Response{
			Status:     status,
			StatusCode: code,
		},
	}
}
