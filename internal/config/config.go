package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/lunjon/http/internal/style"
)

func init() {
	Editor = "vim"
	if v, ok := os.LookupEnv("VISUAL"); ok {
		Editor = v
	} else if v, ok := os.LookupEnv("EDITOR"); ok {
		Editor = v
	}
}

const DefaultConfigString = `
# Valid options and default values
# timeout = "30s"

[aliases] # Section for you URL aliases
# local = http://localhost
`

var (
	Editor         string
	DefaultTimeout = time.Second * 30
)

type Config struct {
	Timeout time.Duration
	Verbose bool
	Fail    bool
	Aliases map[string]string
}

func New() Config {
	return Config{
		Timeout: DefaultTimeout,
		Verbose: false,
		Fail:    false,
		Aliases: make(map[string]string),
	}
}

func (cfg Config) UseVerbose(b bool) Config {
	cfg.Verbose = b
	return cfg
}

func (cfg Config) UseTimeout(t time.Duration) Config {
	cfg.Timeout = t
	return cfg
}

func (cfg Config) UseFail(b bool) Config {
	cfg.Fail = b
	return cfg
}

func (cfg Config) Write(w io.Writer) error {
	encoder := toml.NewEncoder(w)
	return encoder.Encode(cfg)
}

func (cfg Config) String() string {
	var b strings.Builder

	// Root values
	roots := []struct {
		key   string
		value any
	}{
		{"timeout", cfg.Timeout},
	}

	for _, item := range roots {
		key := style.Blue.Render(item.key)
		var val string

		switch value := item.value.(type) {
		case string, time.Duration:
			val = style.None.Render(fmt.Sprintf(`"%s"`, value))
		case bool:
			val = style.BlueB.Render(fmt.Sprint(value))
		default:
			val = style.None.Render(fmt.Sprint(value))
		}

		b.WriteString(fmt.Sprintf("%s = %v\n", key, val))
	}

	if len(cfg.Aliases) > 0 {
		b.WriteString(fmt.Sprintf(
			"\n[%s]\n",
			style.GreenB.Render("aliases"),
		))

		for k, v := range cfg.Aliases {
			line := fmt.Sprintf(`%s = "%s"`, style.Bold.Render(k), v)
			b.WriteString(line + "\n")
		}
	}

	return b.String()
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
		cfg.Timeout.value = DefaultTimeout
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
	Aliases map[string]string
}

func (cfg fileConfig) convert() Config {
	return Config{
		Timeout: cfg.Timeout.value,
		Aliases: cfg.Aliases,
	}
}

type duration struct {
	value time.Duration
}

func (d *duration) UnmarshalText(b []byte) error {
	var err error
	d.value, err = time.ParseDuration(string(b))
	return err
}
