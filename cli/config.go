package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/style"
	"github.com/lunjon/http/internal/util"
)

type ConfigHandler struct {
	configDir  string
	configPath string
	output     io.Writer
}

func newConfigHandler(configPath string, output io.Writer) *ConfigHandler {
	return &ConfigHandler{
		configDir:  path.Dir(configPath),
		configPath: configPath,
		output:     output,
	}
}

func (handler ConfigHandler) list() error {
	cfg, err := config.Load(handler.configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		fmt.Fprintf(
			handler.output,
			"No configuration file found.\nUse %s to create one.\n",
			style.Bold("config init"),
		)
		return nil
	}

	_, err = fmt.Fprintln(handler.output, cfg)
	return err
}

func (handler ConfigHandler) init() error {
	if err := handler.assertDir(); err != nil {
		return err
	}

	_, err := os.Stat(handler.configPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err == nil {
		return nil
	}

	// Create new default config
	f, err := os.Create(handler.configPath)
	if err != nil {
		return err
	}

	_, err = fmt.Print(f, config.DefaultConfigString)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		handler.output,
		`Created new configuration at %s.

Use %s list the configuration, or %s to edit the file.
`,
		style.GreenB(handler.configPath),
		style.Bold("config"),
		style.Bold("config edit"),
	)
	return err
}

func (handler ConfigHandler) edit(editor string) error {
	if err := handler.assertDir(); err != nil {
		return err
	}

	cmd := exec.Command(editor, handler.configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (handler ConfigHandler) assertDir() error {
	exists, isdir, err := util.FileExists(handler.configDir)
	if err != nil {
		return err
	}

	if !exists {
		return os.MkdirAll(handler.configDir, 0700)
	}

	if !isdir {
		return fmt.Errorf(
			"configuration directory %s was expected to be a directory",
			handler.configDir,
		)
	}
	return nil
}
