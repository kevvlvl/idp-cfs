package platform_gp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGoldenPathValidUrl(t *testing.T) {

	gp := GetGoldenPath("https://github.com/kevvlvl/idp-cfs.git", "main", "./gp/my-gp", "1.0.0")

	assert.NotNil(t, gp)
	assert.Equal(t, gp.Tool, GpGithub)
	assert.NotEmptyf(t, gp.URL, "URL is empty")
	assert.NotEmptyf(t, gp.Branch, "Branch is empty")
	assert.NotEmptyf(t, gp.Path, "Path is empty")
	assert.NotEmptyf(t, gp.Tag, "Tag is empty")
}
