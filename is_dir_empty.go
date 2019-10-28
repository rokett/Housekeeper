package main

import (
	"io/ioutil"
)

func isDirEmpty(path string) (bool, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(entries) == 0, nil
}
