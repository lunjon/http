package util

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ignoreDirs = []string{".git", "node_modules", "target"}
)

func ignore(name string) bool {
	return Contains(ignoreDirs, name)
}

func WalkDir(root string) ([]string, error) {
	paths := []string{}
	filepath.WalkDir(root, func(path string, e fs.DirEntry, err error) error {
		if e.IsDir() {
			if ignore(e.Name()) {
				return fs.SkipDir
			}
		} else {
			paths = append(paths, path)
		}

		return nil
	})
	return paths, nil
}

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
