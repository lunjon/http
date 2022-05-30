package complete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComplete(t *testing.T) {
	items := []string{"http://local", "https://example.com", "balloon", "https://golang.org"}
	prefix, matches := Complete("http", items)

	if prefix != "http" {
		t.Errorf("was %s expected %s", prefix, "http")
	}

	if len(matches) != 3 {
		t.Errorf("was %d expected 3", len(matches))
	}
}

func TestCommonPrefix(t *testing.T) {
	tests := []struct {
		items  []string
		prefix string
	}{
		{[]string{"http://local", "https://example.com", "balloon"}, ""},
		{[]string{"http://local", "https://example.com", "https://golang.org"}, "http"},
		{[]string{"https://api", "https://api.example.com", "https://api.org"}, "https://api"},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			actual := commonPrefix(test.items)
			assert.Equal(t, test.prefix, actual)
		})
	}
}
