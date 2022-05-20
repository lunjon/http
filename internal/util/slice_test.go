package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	tests := []struct {
		in  []string
		out []string
		f   func(string) string
	}{
		{[]string{"a", "b"}, []string{"A", "B"}, strings.ToUpper},
		{[]string{"a ", "b\n"}, []string{"a", "b"}, strings.TrimSpace},
		{nil, []string{}, strings.TrimSpace},
		{[]string{}, []string{}, strings.TrimSpace},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			actual := Map(test.in, test.f)
			require.Equal(t, actual, test.out)
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		in  []string
		out []string
		f   func(string) bool
	}{
		{[]string{"a", ""}, []string{"a"}, func(s string) bool { return s != "" }},
		{[]string{"a", "b"}, []string{"a", "b"}, func(s string) bool { return s != "" }},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			actual := Filter(test.in, test.f)
			require.Equal(t, actual, test.out)
		})
	}
}
