package cli

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/format"
	"github.com/spf13/cobra"
)

type FailFunc func(status int)
type runFunc func(*cobra.Command, []string)
type checkRedirectFunc func(*http.Request, []*http.Request) error

type outputs struct {
	logs   io.Writer
	infos  io.Writer
	errors io.Writer
}

const (
	defaultTimeout   = time.Second * 30
	defaultAWSRegion = "eu-west-1"
)

var (
	styler          = format.NewStyler()
	defaultFailFunc = func(int) {}
)

// Build the root command for http and set version.
func Build(version string) (*cobra.Command, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := path.Join(homedir, ".gohttp", "config.yml")
	cfg, err := config.Load(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		cfg = config.New()
	}

	outputs := outputs{
		infos:  os.Stdout,
		logs:   os.Stderr,
		errors: os.Stderr,
	}
	return build(version, cfg, outputs), nil
}

func checkErr(err error, output io.Writer) {
	if err == nil {
		return
	}
	fmt.Fprintf(output, "%s: %v\n", styler.RedB("error"), err)
	os.Exit(1)
}
