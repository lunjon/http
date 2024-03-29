package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/style"
	"github.com/spf13/cobra"
)

type FailFunc func(status int)
type runFunc func(*cobra.Command, []string)

type cliConfig struct {
	logs        io.Writer
	infos       io.Writer
	errors      io.Writer
	configPath  string
	historyPath string
}

func (cfg cliConfig) getAppConfig() (config.Config, error) {
	c, err := config.Load(cfg.configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return c, err
		}
		c = config.New()
	}
	return c, nil
}

const (
	defaultTimeout   = time.Second * 30
	defaultAWSRegion = "eu-west-1"
)

var (
	defaultFailFunc = func(int) {}
)

// Build the root command for http and set version.
func Build(version string) (*cobra.Command, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := path.Join(homedir, ".config", "httpcli")
	configFilepath := path.Join(configDir, "config.toml")
	historyPath := path.Join(configDir, ".history")

	cfg := cliConfig{
		configPath:  configFilepath,
		historyPath: historyPath,
		infos:       os.Stdout,
		logs:        os.Stderr,
		errors:      os.Stderr,
	}
	return build(version, cfg), nil
}

func checkErr(err error, output io.Writer) {
	if err == nil {
		return
	}
	fmt.Fprintf(output, "%s: %v\n", style.RedB.Render("error"), err)
	os.Exit(1)
}
