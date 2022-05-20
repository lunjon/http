package cli

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/lunjon/http/internal/logging"
	"github.com/spf13/cobra"
)

type config struct {
	version        string
	verbose        bool
	trace          bool
	fail           bool
	repeat         int
	defaultHeaders string
	aliasFilepath  string
	output         string
	logs           io.Writer
	infos          io.Writer
	errs           io.Writer
	headerOpt      *HeaderOption
}

func newDefaultConfig(version string) (*config, error) {
	f, err := getAliasFilepath()
	return &config{
		version:        version,
		repeat:         1,
		defaultHeaders: os.Getenv(defaultHeadersEnv),
		aliasFilepath:  f,
		infos:          os.Stdout,
		errs:           os.Stderr,
		logs:           os.Stderr,
		headerOpt:      newHeaderOption(),
	}, err
}

func (c *config) updateFrom(cmd *cobra.Command) error {
	flags := cmd.Flags()
	var err error
	c.verbose, err = flags.GetBool(verboseFlagName)
	if err != nil {
		return err
	}

	c.trace, err = flags.GetBool(traceFlagName)
	if err != nil {
		return err
	}

	c.fail, err = flags.GetBool(failFlagName)
	if err != nil {
		return err
	}

	c.repeat, err = flags.GetInt(repeatFlagName)
	if err != nil {
		return err
	}

	c.output, err = flags.GetString(outputFlagName)
	if err != nil {
		return err
	}

	return nil
}

func (c *config) getLogger() *log.Logger {
	return c.buildLogger(c.verbose)
}

func (c *config) getTraceLogger() *log.Logger {
	return c.buildLogger(c.trace)
}

func (c *config) buildLogger(enabled bool) *log.Logger {
	logger := logging.NewLogger()
	if enabled {
		logger.SetOutput(c.logs)
	} else {
		logger.SetOutput(ioutil.Discard)
	}
	return logger
}
