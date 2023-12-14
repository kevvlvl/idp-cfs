package util

import "os"

func RemoveAllDir(path string) error {
	return os.RemoveAll(path)
}
