package global

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"path/filepath"
)

// CreateFolder creates a folder if it does not exist. Returns error if the folder exists
func CreateFolder(dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			log.Error().Msgf("CreateFolder() - Failed to create directory: %v", err)
			return err
		}
	} else {
		e := errors.New("CreateFolder() - directory exists! Please make sure the dir does not exist")
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
					log.Error().Msgf("CopyFilesDeep() - Failed to close the src file %s: %v", srcFile.Name(), err)
				}
			}(srcFile)

			destFilePath := filepath.Join(destDir, info.Name())
			destFile, _ := os.Create(destFilePath)
			defer func(destFile *os.File) {
				err := destFile.Close()
				if err != nil {
					log.Error().Msgf("CopyFilesDeep() - Failed to close the src file %s: %v", srcFile.Name(), err)
				}
			}(destFile)

			_, err := io.Copy(destFile, srcFile)
			if err != nil {
				log.Error().Msgf("CopyFilesDeep() - Failed to copy the file from gp path to the new code path: %v", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Error().Msgf("CopyFilesDeep() - Failed to walk the directory %s: %v", srcDir, err)
		return err
	}

	return nil
}

func ExtractTgz(gzStream io.Reader, dst string) {

	uncompressed, err := gzip.NewReader(gzStream)

	if err != nil {
		log.Error().Msgf("ExtractTgz() - failed to Read gzStream: %s", err)
	}

	tarReader := tar.NewReader(uncompressed)

	for {

		h, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error().Msgf("ExtractTgz() - Reading tar failed: %s", err)
		}

		dstPath := path.Join(dst, h.Name)

		switch h.Typeflag {

		case tar.TypeDir:

			err := CreateFolder(dstPath)
			if err != nil {
				break
			}
		case tar.TypeReg:
			f, err := os.Create(dstPath)

			if err != nil {
				log.Error().Msgf("ExtractTgz() - Failed to create file %s: %s", dstPath, err)
			}

			if _, err := io.Copy(f, tarReader); err != nil {
				log.Error().Msgf("ExtractTgz() - File copy failed: %s", err)
			}

			err = f.Close()
			if err != nil {
				log.Error().Msgf("ExtractTgz() - Failed to close the file %s: %s", dstPath, err)
			}
		default:
			log.Error().Msgf("ExtractTgz() - Unexpected tar type: %v. Name: %s", h.Typeflag, dstPath)
		}
	}
}
