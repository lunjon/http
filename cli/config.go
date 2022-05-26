package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/lunjon/http/internal/config"
)

var ()

const ()

type ConfigHandler struct {
	configPath string
	output     io.Writer
}

func newConfigHandler(configPath string, output io.Writer) *ConfigHandler {
	return &ConfigHandler{
		configPath: configPath,
		output:     output,
	}
}

func (handler *ConfigHandler) list() error {
	cfg, err := config.Load(handler.configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		fmt.Fprintf(
			handler.output,
			`No configuration file found.
Use %s to create one.`+"\n",
			styler.WhiteB("config init"),
		)
		return nil
	}

	fmt.Fprintln(handler.output, cfg.String())
	return nil
}

func (handler *ConfigHandler) init() error {
	panic("not implemented")
}

func (handler *ConfigHandler) edit() error {
	panic("not implemented")
}
