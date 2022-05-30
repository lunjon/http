package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	tests := []struct {
		path   string
		exists bool
		isdir  bool
	}{
		{"testdata/test.json", true, false},
		{"testdata", true, true},
		{"unknown-file", false, false},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			exists, isdir, err := FileExists(test.path)
			assert.NoError(t, err)
			assert.Equal(t, test.exists, exists)
			assert.Equal(t, test.isdir, isdir)
		})
	}

}
