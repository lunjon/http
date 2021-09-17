package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDefault(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"get", serverURL})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestGetSilent(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"get", serverURL, "--silent"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestGetBrief(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"get", serverURL, "--brief"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestPostDefault(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"post", serverURL, "--body", `{"data": "string"}`})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestWithVerbose(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"post", serverURL, "--body", `{"data": "string"}`, "--verbose"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestDeleteWithFail(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"delete", serverURL, "--fail"})

	err := fixture.root.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestMissingURLParam(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"get"})

	err := fixture.root.Execute()
	require.Error(t, err)
}

func TestUnknownCommand(t *testing.T) {
	fixture := setup(t)
	fixture.root.SetArgs([]string{"unknown"})

	err := fixture.root.Execute()
	require.Error(t, err)
}
