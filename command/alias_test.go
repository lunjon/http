package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAliasListEmpty(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"alias"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestAliasList(t *testing.T) {
	fixture := setup(t)
	fixture.handler.setAlias("local", "http://localhost")
	fixture.root.SetArgs([]string{"alias"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestAliasSet(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"alias", "ss", "http://localhost"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())

	aliases, err := fixture.handler.readAliasFile()
	require.NoError(t, err)
	require.NotEmpty(t, aliases)
}

func TestAliasInvalidParams(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"alias", "odd"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.errors)
}
