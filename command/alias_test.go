package command

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type aliasManagerMock struct {
	aliases map[string]string
}

func newAliasManagerMock() *aliasManagerMock {
	return &aliasManagerMock{make(map[string]string)}
}

func (m *aliasManagerMock) set(name, value string) {
	m.aliases[name] = value
}

func (m *aliasManagerMock) Load() (map[string]string, error) {
	return m.aliases, nil
}

func (m *aliasManagerMock) Save(aliases map[string]string) error {
	m.aliases = aliases
	return nil
}

type aliasTestFixture struct {
	handler *AliasHandler
	infos   *strings.Builder
	errors  *strings.Builder
}

func setupAliasHandlerTest(t *testing.T) *aliasTestFixture {
	infos := &strings.Builder{}
	errors := &strings.Builder{}

	m := newAliasManagerMock()

	h := &AliasHandler{
		manager: m,
		infos:   infos,
		errors:  errors,
	}

	return &aliasTestFixture{
		handler: h,
		infos:   infos,
		errors:  errors,
	}
}

func TestAliasList(t *testing.T) {
	fixture := setupAliasHandlerTest(t)
	err := fixture.handler.setAlias("local", "http://localhost")
	require.NoError(t, err)

	err = fixture.handler.listAlias()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestAliasRemove(t *testing.T) {
	fixture := setupAliasHandlerTest(t)
	err := fixture.handler.setAlias("local", "http://localhost")
	require.NoError(t, err)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"local", false},
		{"unknown", true},
		{"", true},
		{"a!", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := fixture.handler.removeAlias(test.name)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAliasSet(t *testing.T) {
	fixture := setupAliasHandlerTest(t)
	err := fixture.handler.setAlias("local", "http://localhost")
	require.NoError(t, err)
	require.Empty(t, fixture.errors.String())
}

func TestAliasSetInvalid(t *testing.T) {
	names := []string{
		"",
		"1",
		"^",
		"a#",
		"yEs!",
	}

	fixture := setupAliasHandlerTest(t)
	for _, name := range names {
		err := fixture.handler.setAlias(name, "http://localhost")
		require.Error(t, err)
	}
}

func setupAliasManagerTest(t *testing.T, aliases map[string]string) *fileAliasManager {
	if aliases != nil {
		b, err := json.Marshal(aliases)
		if err != nil {
			t.Fatalf("error on setup: %s", err)
		}
		err = os.WriteFile(testAliasFilepath, b, 0600)
		if err != nil {
			t.Fatalf("error on setup: %s", err)
		}
	}
	t.Cleanup(func() {
		os.Remove(testAliasFilepath)
	})
	return newAliasLoader(testAliasFilepath)
}

func TestFileAliasaManagerLoad(t *testing.T) {
	m := setupAliasManagerTest(t, nil)

	aliases, err := m.Load()
	require.NoError(t, err)
	assert.Empty(t, aliases)
}

func TestFileAliasaManagerLoadWithData(t *testing.T) {
	m := setupAliasManagerTest(t, map[string]string{
		"test": "http://localhost/path",
	})

	aliases, err := m.Load()
	require.NoError(t, err)
	assert.Len(t, aliases, 1)
}

func TestFileAliasaManagerSave(t *testing.T) {
	m := setupAliasManagerTest(t, nil)
	err := m.Save(map[string]string{
		"test": "http://localhost/path",
	})
	require.NoError(t, err)
}
