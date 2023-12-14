package client_github

import (
	"errors"
	"github.com/google/go-github/v56/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidateApiResponseWithNilError_ThenNil(t *testing.T) {

	resultErr := validateApiResponse(respClientError(), nil, "Client error")

	assert.Nil(t, resultErr, "validateApiResponse returned non-nil error for a nil error input")
}

func TestValidateApiResponseWithClientError_ThenErrorHttp4xx(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := validateApiResponse(respClientError(), err, msg)

	assert.NotNil(t, resultErr, "validateApiResponse returned nil error for a HTTP4xx error")
	assert.Contains(t, resultErr.Error(), "HTTP4xx - ")
}

func TestValidateApiResponseWithClientError_ThenErrorHttp404(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := validateApiResponse(respClientNotFoundError(), err, msg)

	assert.NotNil(t, resultErr, "validateApiResponse returned nil error for a HTTP404 error")
	assert.Contains(t, resultErr.Error(), "HTTP404 - ")
}

func TestValidateApiResponseWithClientError_ThenErrorHttp5xx(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := validateApiResponse(respServerError(), err, msg)

	assert.NotNil(t, resultErr, "validateApiResponse returned nil error for a HTTP5xx error")
	assert.Contains(t, resultErr.Error(), "HTTP5xx - ")
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
