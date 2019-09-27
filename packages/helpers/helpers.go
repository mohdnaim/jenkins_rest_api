package helpers

import (
	"os"
	"path/filepath"
)

// StringInSlice ...
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetFilenamesRecursively ...
func GetFilenamesRecursively(path string) []string {
	var files []string
	err := filepath.Walk("xml", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
