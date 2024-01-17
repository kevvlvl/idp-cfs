package global

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func errorTestNotNil() (error, string) {

	errorMsg := "failed to query remote git system"
	return errors.New(errorMsg), errorMsg
}

func getResponse(status string, code int) *http.Response {
	return &http.Response{
		Status:     status,
		StatusCode: code,
	}
}

func respClientError() *http.Response {
	return getResponse("Bad Request", 400)
}

func respClientNotFoundError() *http.Response {
	return getResponse("Not Found", 404)
}

func respServerError() *http.Response {
	return getResponse("Internal Server Error", 500)
}

func TestValidateApiResponseWithNilError_ThenNil(t *testing.T) {

	resultErr := ValidateApiResponse(respClientError(), nil, "Client error")

	assert.Nil(t, resultErr)
}

func TestValidateApiResponseWithClientError_ThenErrorHttp4xx(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := ValidateApiResponse(respClientError(), err, msg)

	assert.NotNil(t, resultErr)
	assert.Contains(t, resultErr.Error(), "HTTP4xx - ")
}

func TestValidateApiResponseWithClientError_ThenErrorHttp404(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := ValidateApiResponse(respClientNotFoundError(), err, msg)

	assert.NotNil(t, resultErr)
	assert.Contains(t, resultErr.Error(), "HTTP404 - ")
}

func TestValidateApiResponseWithClientError_ThenErrorHttp5xx(t *testing.T) {

	err, msg := errorTestNotNil()
	resultErr := ValidateApiResponse(respServerError(), err, msg)

	assert.NotNil(t, resultErr)
	assert.Contains(t, resultErr.Error(), "HTTP5xx - ")
}
