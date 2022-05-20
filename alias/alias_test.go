package alias

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/lunjon/http/format"
	"github.com/lunjon/http/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testdir           = "test-http"
	testAliasFilepath = path.Join(testdir, "aliases.json")
)

func TestMain(m *testing.M) {
	if _, err := os.Stat(testdir); os.IsNotExist(err) {
		err := os.MkdirAll(testdir, 0700)
		if err != nil {
			panic(err)
		}
	}

	status := m.Run()
	os.RemoveAll(testdir)
	os.Exit(status)
}

type aliasTestFixture struct {
	handler *Handler
	infos   *strings.Builder
	errors  *strings.Builder
}

func setupAliasHandlerTest(t *testing.T) *aliasTestFixture {
	infos := &strings.Builder{}
	errors := &strings.Builder{}

	m := mock.NewManagerMock()
	h := NewHandler(m, format.NewStyler(), infos, errors)

	return &aliasTestFixture{
		handler: h,
		infos:   infos,
		errors:  errors,
	}
}

func TestAliasList(t *testing.T) {
	fixture := setupAliasHandlerTest(t)
	err := fixture.handler.Set("local", "http://localhost")
	require.NoError(t, err)

	err = fixture.handler.List(true)
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestAliasRemove(t *testing.T) {
	fixture := setupAliasHandlerTest(t)
	err := fixture.handler.Set("local", "http://localhost")
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
			err := fixture.handler.Remove(test.name)
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
	err := fixture.handler.Set("local", "http://localhost")
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
		err := fixture.handler.Set(name, "http://localhost")
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
	return NewManager(testAliasFilepath)
}

func TestFileAliasManagerLoad(t *testing.T) {
	m := setupAliasManagerTest(t, nil)

	aliases, err := m.Load()
	require.NoError(t, err)
	assert.Empty(t, aliases)
}

func TestFileAliasManagerLoadWithData(t *testing.T) {
	m := setupAliasManagerTest(t, map[string]string{
		"test": "http://localhost/path",
	})

	aliases, err := m.Load()
	require.NoError(t, err)
	assert.Len(t, aliases, 1)
}

func TestFileAliasManagerSave(t *testing.T) {
	m := setupAliasManagerTest(t, nil)
	err := m.Save(map[string]string{
		"test": "http://localhost/path",
	})
	require.NoError(t, err)
}
