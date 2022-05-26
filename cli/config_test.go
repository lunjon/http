package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigList(t *testing.T) {
	output := &strings.Builder{}
	configPath := testConfigPath
	handler := newConfigHandler(configPath, output)

	err := handler.list()
	assert.NoError(t, err)
	assert.NotEmpty(t, output.String())
}
