package main

import (
	"errors"
	"os"
)

func optStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optFloat64(n float64) *float64 {
	if n == 0 {
		return nil
	}
	return &n
}

func openDir(path string) (dir *os.File, err error) {
	// Open the directory
	dir, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	// Fail if the file isn't a directory.
	dirInfo, err := dir.Stat()
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, errors.New("openDir: not a directory")
	}

	return dir, err
}
