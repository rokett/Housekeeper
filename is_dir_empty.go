package main

import (
	"fmt"
	"io/ioutil"
)

func isDirEmpty(path string) (bool, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	// In order to tell whether a directory is empty, we need to check any sub directories as well.
	for _, dir := range entries {
		if dir.IsDir() == false {
			return false, nil
		}

		subdir := fmt.Sprintf("%s\\%s", path, dir.Name())

		empty, err := isDirEmpty(subdir)
		if err != nil {
			return false, err
		}

		if empty == false {
			return false, nil
		}
	}

	return true, nil
}
