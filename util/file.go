package util

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
)

// RemoveAllDir removes the path (dir) and subdirectories and files
func RemoveAllDir(path string) error {
	return os.RemoveAll(path)
}

// CopyFilesDeep copies all files from srcDir to destDir recursively
func CopyFilesDeep(srcDir string, destDir string) error {

	err := filepath.Walk(srcDir, func(file string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			srcFile, _ := os.Open(file)
			defer func(srcFile *os.File) {
				err := srcFile.Close()
				if err != nil {
					log.Error().Msgf("Failed to close the src file %s: %v", srcFile.Name(), err)
				}
			}(srcFile)

			destFilePath := filepath.Join(destDir, info.Name())
			destFile, _ := os.Create(destFilePath)
			defer func(destFile *os.File) {
				err := destFile.Close()
				if err != nil {
					log.Error().Msgf("Failed to close the src file %s: %v", srcFile.Name(), err)
				}
			}(destFile)

			_, err := io.Copy(destFile, srcFile)
			if err != nil {
				log.Error().Msgf("Failed to copy the file from gp path to the new code path: %v", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Error().Msgf("Failed to walk the directory %s: %v", srcDir, err)
		return err
	}

	return nil
}
