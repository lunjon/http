package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const defaultConfig = `
# Valid options and default values
# timeout = "30s"
# verbose = false
# fail = false

[aliases] # Section for you URL aliases
# local = http://localhost
`

type Config struct {
	Timeout time.Duration
	Verbose bool
	Fail    bool
	Aliases map[string]string
}

func New() Config {
	return Config{
		Verbose: false,
		Fail:    false,
		Aliases: make(map[string]string),
	}
}

func (cfg Config) UseVerbose(b bool) Config {
	cfg.Verbose = b
	return cfg
}

func (cfg Config) UseFail(b bool) Config {
	cfg.Fail = b
	return cfg
}

// ReadTOML loads the Config from a TOML formatted byte slice.
func ReadTOML(data []byte) (Config, error) {
	var cfg fileConfig
	err := toml.Unmarshal(data, &cfg)
	if err != nil {
		return New(), err
	}

	// Correction of zero values
	if cfg.Timeout.value == 0 {
		cfg.Timeout.value = time.Second * 30
	}

	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
	}
	return cfg.convert(), err
}

func Load(filepath string) (Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return New(), err
	}
	return ReadTOML(data)
}

type fileConfig struct {
	Timeout duration
	Verbose bool
	Fail    bool
	Aliases map[string]string
}

func (cfg fileConfig) convert() Config {
	return Config{
		Timeout: cfg.Timeout.value,
		Verbose: cfg.Verbose,
		Fail:    cfg.Fail,
		Aliases: cfg.Aliases,
	}
}

type duration struct {
	value time.Duration
}

func (d *duration) UnmarshalText(b []byte) error {
	fmt.Println("YES")
	fmt.Println(string(b))
	var err error
	d.value, err = time.ParseDuration(string(b))
	return err
}
