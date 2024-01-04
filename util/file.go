package util

import (
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
)

// CreateFolder creates a folder if it does not exist. Returns error if the folder exists
func CreateFolder(dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			log.Error().Msgf("Failed to create directory: %v", err)
			return err
		}
	} else {
		e := errors.New("directory exists! Please make sure the dir does not exist")
		log.Error().Msg(e.Error())
		return e
	}

	return nil
}

// CopyFilesDeep copies all files from srcDir to destDir recursively
func CopyFilesDeep(srcDir, destDir string) error {

	err := filepath.Walk(srcDir, func(file string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

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
