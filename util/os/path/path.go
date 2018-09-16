// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package path defines some utils about path
package path

import (
	"os"
	"path/filepath"
)

// GetDirPath returns current dir path
func GetCurrentDirPath() string {
	dirPath := filepath.Dir(os.Args[0])
	return dirPath
}

// GetPath returns abs path of the file
func GetPath(file string) (string, error) {
	if filepath.IsAbs(file) {
		return file, nil
	} else {
		path := GetCurrentDirPath()
		return path + "/" + file, nil
	}
}
