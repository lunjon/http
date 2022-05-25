package config

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const defaultConfig = `
# Valid root options and default values
# timeout = "30s"
# verbose = false
# trace = false
# fail = false
# repeat = 0

[aliases] # Section for you URL aliases
# local = http://localhost
`

type Config struct {
	Timeout time.Duration
	Verbose bool
	Trace   bool
	Fail    bool
	Repeat  int
	Aliases map[string]string
}

func New() Config {
	return Config{
		Verbose: false,
		Trace:   false,
		Fail:    false,
		Repeat:  1,
		Aliases: make(map[string]string),
	}
}

func (cfg Config) UseVerbose(b bool) Config {
	cfg.Verbose = b
	return cfg
}

func (cfg Config) UseTrace(b bool) Config {
	cfg.Trace = b
	return cfg
}

func (cfg Config) UseFail(b bool) Config {
	cfg.Fail = b
	return cfg
}

func (cfg Config) UseRepeat(n int) Config {
	cfg.Repeat = n
	return cfg
}

// ReadTOML loads the Config from a TOML formatted byte slice.
func ReadTOML(data []byte) (Config, error) {
	var cfg Config
	err := toml.Unmarshal(data, &cfg)

	// Correction of zero values
	if cfg.Timeout == 0 {
		cfg.Timeout = time.Second * 30
	}

	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
	}

	return cfg, err
}

func Load(filepath string) (Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}
	return ReadTOML(data)
}
