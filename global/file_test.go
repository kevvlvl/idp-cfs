package global

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFolder_New_NoError(t *testing.T) {

	testFolder := "/tmp/idp-cfs-unittest-util"
	err := CreateFolder(testFolder)

	assert.Nil(t, err)

	err = os.RemoveAll(testFolder)
	assert.Nil(t, err)
}

func TestCreateFolder_Exists_Error(t *testing.T) {

	testFolder := "/tmp"
	err := CreateFolder(testFolder)

	assert.NotNil(t, err)
	assert.Equal(t, "directory exists! Please make sure the dir does not exist", err.Error())
}

func TestCopyFilesDeep_SrcValid_DstValid_NoError(t *testing.T) {

	// Prepare test folders

	srcDir := "/tmp/idp-cfs-copy-src"
	srcFile := filepath.Join(srcDir, "/test-file.txt")
	dstDir := "/tmp/idp-cfs-copy-dst"

	err := os.Mkdir(srcDir, 0755)
	assert.Nil(t, err)

	err = os.Mkdir(dstDir, 0755)
	assert.Nil(t, err)

	// Create sample file

	newTestFile, err := os.Create(srcFile)
	assert.Nil(t, err)
	assert.NotNil(t, newTestFile)

	err = CopyFilesDeep(srcDir, dstDir)
	assert.Nil(t, err)

	// Verify dst folder

	_, err = os.Stat(srcFile)
	assert.Nil(t, err)

	// Cleanup!
	err = os.RemoveAll(srcDir)
	assert.Nil(t, err)

	err = os.RemoveAll(dstDir)
	assert.Nil(t, err)
}

func TestCopyFilesDeep_SrcNotExist_Error(t *testing.T) {

	srcDir := "/tmp/idp-cfs-some-nonexistant-folder"
	dstDir := "/tmp/idp-cfs-unitest-util-copy"

	err := CopyFilesDeep(srcDir, dstDir)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}
