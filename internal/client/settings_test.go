package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCertOptionsFromEmpty(t *testing.T) {
	opts, err := CertOptionsFrom("", "")
	require.NoError(t, err)
	require.True(t, opts.IsNone())
}

func TestCertOptionsFromOnlyKey(t *testing.T) {
	_, err := CertOptionsFrom("", "key")
	require.Error(t, err)
}
