package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	cfg, err := ReadTOML([]byte(defaultConfig))
	assert.NoError(t, err)

	assert.Equal(t, time.Second*30, cfg.Timeout)
	assert.False(t, cfg.Verbose)
	assert.False(t, cfg.Fail)
	assert.Len(t, cfg.Aliases, 0)
}

func TestTimeout(t *testing.T) {
	s := `timeout = "13s"`
	cfg, err := ReadTOML([]byte(s))
	assert.NoError(t, err)
	assert.Equal(t, time.Second*13, cfg.Timeout)
}

func TestAliases(t *testing.T) {
	s := `[aliases]
local = "https://localhost/path"`
	cfg, err := ReadTOML([]byte(s))
	assert.NoError(t, err)
	assert.Len(t, cfg.Aliases, 1)
	assert.Equal(t, "https://localhost/path", cfg.Aliases["local"])
}

func TestString(t *testing.T) {
	cfg, err := ReadTOML([]byte(defaultConfig))
	assert.NoError(t, err)
	assert.NotZero(t, cfg.String())
}
