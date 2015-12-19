package io

import (
	"os"
	"path/filepath"
)

func DirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	return false
}

func DirChildren(root string, predicate func(string, os.FileInfo, error) bool) []string {
	r := make([]string, 0)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if predicate(path, info, err) {
			r = append(r, path)
		}
		return nil
	}
	filepath.Walk(root, walkFn)
	return r
}
