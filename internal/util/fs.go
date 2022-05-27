package util

import (
	"errors"
	"os"
	"os/exec"
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

// OpenEditor edits a temporary file and returns the content.
func OpenEditor(editor string) ([]byte, error) {
	f, err := os.CreateTemp(".", "body")
	if err != nil {
		return nil, err
	}

	name := f.Name()
	if err = f.Close(); err != nil {
		return nil, err
	}
	defer func() {
		os.Remove(name)
	}()

	cmd := exec.Command(editor, name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(name)
}
