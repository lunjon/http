package command

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type aliasTestFixture struct {
	handler *AliasHandler
	infos   *strings.Builder
	errors  *strings.Builder
}

func setupAliasTest(t *testing.T) *aliasTestFixture {
	infos := &strings.Builder{}
	errors := &strings.Builder{}

	h := &AliasHandler{
		aliasFilepath: testAliasFilepath,
		infos:         infos,
		errors:        errors,
	}

	return &aliasTestFixture{
		handler: h,
		infos:   infos,
		errors:  errors,
	}
}

func TestAliasList(t *testing.T) {
	fixture := setupAliasTest(t)
	err := fixture.handler.setAlias("local", "http://localhost")
	require.NoError(t, err)

	err = fixture.handler.listAlias()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestAliasRemove(t *testing.T) {
	fixture := setupAliasTest(t)
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
	fixture := setupAliasTest(t)
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

	fixture := setupAliasTest(t)
	for _, name := range names {
		err := fixture.handler.setAlias(name, "http://localhost")
		require.Error(t, err)
	}
}
