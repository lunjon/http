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
}
