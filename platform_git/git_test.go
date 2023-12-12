package platform_git

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCodeGitlabNotImplemented_ThenNil(t *testing.T) {
	assert.Nil(t, GetCode("gitlab", "/tmp/unit-test"))
}

func TestGetCodeGiteaNotImplemented_ThenNil(t *testing.T) {
	assert.Nil(t, GetCode("gitea", "/tmp/unit-test"))
}
