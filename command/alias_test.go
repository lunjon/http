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

func TestAliasListEmpty(t *testing.T) {
	fixture := setupAliasTest(t)
	err := fixture.handler.listAlias()

	require.NoError(t, err)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
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

func TestAliasSet(t *testing.T) {
	fixture := setupAliasTest(t)
	err := fixture.handler.setAlias("local", "http://localhost")
	require.NoError(t, err)
	require.Empty(t, fixture.errors.String())
}
