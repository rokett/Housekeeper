package main

import (
	"io"
	"os"
)

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read a single file and if it is EOF, the directory is empty
	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
