package tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const Separator = string(filepath.Separator)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func IsDirectory(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return stat.IsDir(), nil
}

func PathWithSeparator(path string) string {
	if strings.HasSuffix(path, Separator) {
		return path
	}

	return path + Separator
}

func PathWithoutSeparator(path string) string {
	if strings.HasSuffix(path, Separator) {
		return PathWithoutSeparator(path[:len(path)-1])
	}

	return path
}

func CreateDirectoryIfNotExists(path string) (string, error) {
	var exists bool
	var err error

	if exists, err = FileExists(path); err != nil {
		return "", err
	}

	if !exists {
		// create directory
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return "", err
		}

		return PathWithSeparator(path), nil
	}

	var isDir bool
	if isDir, err = IsDirectory(path); err != nil {
		return "", err
	}

	if !isDir {
		return "", errors.New(fmt.Sprintf("file '%v' is not a directory", path))
	}

	return PathWithSeparator(path), nil
}
