package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaber(t *testing.T) {
	taber := NewTaber("")
	taber.Writef("%s:\n", "header")
	taber.WriteLine("one", "two")

	s := taber.String()
	require.NotEmpty(t, s)
	lines := strings.Split(s, "\n")
	require.Len(t, lines, 3)
	require.Equal(t, "header:", lines[0])
}
