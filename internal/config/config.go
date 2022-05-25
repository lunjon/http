package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
)

const defaultConfig = `
# Valid options and default values
# timeout = "30s"
# verbose = false
# fail = false

[aliases] # Section for you URL aliases
# local = http://localhost
`

var (
	noStyle       = lipgloss.NewStyle()
	nameStyle     = lipgloss.NewStyle().Bold(true)
	boolStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))
	durationStyle = lipgloss.NewStyle()
	sectionStyle  = nameStyle.Copy().Foreground(lipgloss.Color("2"))
)

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

func (cfg Config) String() string {
	var b strings.Builder

	// Root values
	roots := []struct {
		key   string
		value any
	}{
		{"timeout", cfg.Timeout},
		{"verbose", cfg.Verbose},
		{"fail", cfg.Fail},
	}

	for _, item := range roots {
		key := nameStyle.Render(item.key)
		var val string
		switch value := item.value.(type) {
		case string:
			val = noStyle.Render(fmt.Sprintf(`"%s"`, value))
		case bool:
			val = boolStyle.Render(fmt.Sprint(value))
		case time.Duration:
			val = durationStyle.Render(fmt.Sprint(value))

		}
		b.WriteString(fmt.Sprintf("%s = %v\n", key, val))
	}

	if len(cfg.Aliases) > 0 {
		b.WriteString(fmt.Sprintf(
			"\n[%s]\n",
			sectionStyle.Render("aliases"),
		))

		for k, v := range cfg.Aliases {
			line := fmt.Sprintf(` %s = "%s"`, nameStyle.Render(k), v)
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
