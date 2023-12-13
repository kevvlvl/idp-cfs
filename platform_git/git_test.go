package platform_git

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetCodeGitlabNotImplemented_ThenNil(t *testing.T) {
	assert.Nil(t, GetCode("gitlab", "/tmp/unit-test"))
}

func TestGetCodeGiteaNotImplemented_ThenNil(t *testing.T) {
	assert.Nil(t, GetCode("gitea", "/tmp/unit-test"))
}

func TestDeleteCodePath_ValidPath_NoErrors(t *testing.T) {

	testCheckoutPath := "/tmp/idp-cfs-code-unit-test"

	gc := GitCode{
		CodeClonePath: testCheckoutPath,
	}

	f, err := os.Create(testCheckoutPath)
	if err != nil {
		t.Fatalf("Failed to create directory for unit test: %v", err)
	}

	t.Logf("Created temp folder %s for unit test", f.Name())

	err = gc.DeleteCodePath()
	assert.Nil(t, err)
}
