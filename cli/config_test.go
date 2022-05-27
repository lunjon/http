package cli

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type configTest struct {
	configPath string
	output     *strings.Builder
	h          *ConfigHandler
}

func setupConfigTest(t *testing.T) *configTest {
	output := &strings.Builder{}
	configPath := path.Join(testdir, "config-test.toml")
	t.Cleanup(func() {
		_ = os.Remove(configPath)
	})

	h := newConfigHandler(configPath, output)
	return &configTest{
		configPath: configPath,
		output:     output,
		h:          h,
	}
}

func TestConfigList(t *testing.T) {
	test := setupConfigTest(t)

	err := test.h.list()
	assert.NoError(t, err)
	assert.NotEmpty(t, test.output.String())
}

func TestConfigInit(t *testing.T) {
	test := setupConfigTest(t)
	_ = os.Remove(test.configPath)

	err := test.h.init()
	assert.NoError(t, err)
	assert.Contains(t, test.output.String(), "Created new")
}
