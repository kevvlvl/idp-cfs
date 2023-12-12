package platform_gp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
