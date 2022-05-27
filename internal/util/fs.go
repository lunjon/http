package util

import (
	"errors"
	"os"
)

func FileExists(filepath string) (exists, isdir bool, err error) {
	var stat os.FileInfo
	stat, err = os.Stat(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
		return
	}

	exists = true
	isdir = stat.IsDir()
	return
}
