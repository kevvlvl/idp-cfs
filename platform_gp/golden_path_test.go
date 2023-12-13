package platform_gp

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDeleteClonePathDir_ValidPath_NoErrors(t *testing.T) {

	testCheckoutPath := "/tmp/idp-cfs-gp-unit-test"

	gp := GoldenPath{
		GpCheckoutPath: testCheckoutPath,
	}

	f, err := os.Create(testCheckoutPath)
	if err != nil {
		t.Fatalf("Failed to create directory for unit test: %v", err)
	}

	t.Logf("Created temp folder %s for unit test", f.Name())

	err = gp.DeleteClonePathDir()
	assert.Nil(t, err)
}

func TestGetGoldenPathValidUrl(t *testing.T) {

	gp := GetGoldenPath("https://github.com/kevvlvl/idp-cfs.git", "main", "./gp/my-gp", "1.0.0", DefaultCheckoutPath)

	assert.NotNil(t, gp)
	assert.Equal(t, GpGithub, gp.Tool)
	assert.NotEmptyf(t, gp.URL, "URL is empty")
	assert.NotEmptyf(t, gp.Branch, "Branch is empty")
	assert.NotEmptyf(t, gp.Path, "Path is empty")
	assert.NotEmptyf(t, gp.Tag, "Tag is empty")
	assert.NotEmptyf(t, gp.GpCheckoutPath, "Checkout path is empty")
}

func TestFailedCloneGpError(t *testing.T) {
	e := failedCloneGpError()

	assert.NotNil(t, e)
	assert.Equal(t, "failed to clone the GoldenPath", e.Error())
}
