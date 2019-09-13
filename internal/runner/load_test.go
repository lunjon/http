package runner_test

import (
	"testing"

	"github.com/lunjon/httpreq/internal/runner"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		// JSON
		{"minimal json spec", "testdata/minimal.json", false},
		{"post", "testdata/post.json", false},
		{"headers", "testdata/headers.json", false},
		{"ID with whitespace", "testdata/id_whitespace.json", true},
		{"wrong method json spec", "testdata/wrong_method.json", true},
		{"headers list", "testdata/headers_list.json", true},
		{"invalid URL json spec", "testdata/invalid_url.json", true},
		{"missing ID", "testdata/missing_id.json", true},
		{"missing URL", "testdata/missing_url.json", true},
		{"post, missing body", "testdata/post_missing_body.json", true},
		// YAML
		{"minimal yaml spec", "testdata/minimal.yaml", false},
		{"minimal yml spec", "testdata/minimal.yml", false},
		{"headers 1", "testdata/headers.yaml", false},
		{"headers 2", "testdata/headers.yml", false},
		{"post", "testdata/post.yml", false},
		{"empty", "testdata/empty.yaml", true},
		{"post, missing body", "testdata/post_missing_body.yml", true},
		{"headers, list", "testdata/headers_list.yaml", true},
		// Other
		{"unknown file", "unknown", true},
		{"invalid file extensions", "unknown.txt", true},
		{"unique names", "testdata/unique_ids.yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, err := runner.Load(tt.filepath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && runner == nil {
				t.Errorf("Load() returned nil")
			}
		})
	}
}
