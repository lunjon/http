package command

import (
	"io"
	"os"
)

type config struct {
	version        string
	fail           bool
	repeat         int
	defaultHeaders string
	aliasFilepath  string
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

func (c *config) setRepeat(val int) *config {
	c.repeat = val
	return c
}

func (c *config) setFail(val bool) *config {
	c.fail = val
	return c
}
