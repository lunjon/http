package runner_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/lunjon/httpreq/internal/runner"
)

func TestLoadGoodJSON(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"minimal json spec", "testdata/json/minimal.json", false},
		{"post", "testdata/json/post.json", false},
		{"headers", "testdata/json/headers.json", false},
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

func TestLoadBadJSON(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"ID with whitespace", "testdata/json/id_whitespace.json", true},
		{"wrong method json spec", "testdata/json/wrong_method.json", true},
		{"headers list", "testdata/json/headers_list.json", true},
		{"invalid URL json spec", "testdata/json/invalid_url.json", true},
		{"missing ID", "testdata/json/missing_id.json", true},
		{"missing URL", "testdata/json/missing_url.json", true},
		{"post, missing body", "testdata/json/post_missing_body.json", true},

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

func TestLoadGoodYAML(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"minimal yaml spec", "testdata/yaml/minimal.yaml", false},
		{"minimal yml spec", "testdata/yaml/minimal.yml", false},
		{"headers 1", "testdata/yaml/headers.yaml", false},
		{"headers 2", "testdata/yaml/headers.yml", false},
		{"post", "testdata/yaml/post.yml", false},

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

func TestLoadBadYAML(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"empty", "testdata/yaml/empty.yaml", true},
		{"post, missing body", "testdata/yaml/post_missing_body.yml", true},
		{"headers, list", "testdata/yaml/headers_list.yaml", true},
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

func TestLoadMisc(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{"unknown file", "unknown", true},
		{"invalid file extensions", "unknown.txt", true},
		{"unique names", "testdata/yaml/unique_ids.yaml", true},
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

func TestLoadDefaultHeaders(t *testing.T) {
	rn, err := runner.Load("testdata/json/default_headers.json")
	assert.NoError(t, err)
	assert.NotNil(t, rn.Spec.Headers)
	r := rn.Spec.Requests[0]

	assert.Contains(t, r.Headers, "name")
	assert.Contains(t, r.Headers, "token")

	assert.Equal(t, r.Headers["name"], "override")
	assert.Equal(t, r.Headers["token"], "secret")
}

func TestLoadAWSSigv4Defaults(t *testing.T) {
	rn, err := runner.Load("testdata/yaml/aws-sigv4-bool.yml")
	assert.NoError(t, err)
	assert.NotNil(t, rn)

	r := rn.Spec.Requests[0]
	assert.NotNil(t, r.AWS)
	aws := r.GetAWSSign()
	assert.NotNil(t, aws)
	assert.Equal(t, "eu-west-1", aws.Region)
	assert.Empty(t, aws.Profile)
}


func TestLoadAWSSigv4RegionOnly(t *testing.T) {
	rn, err := runner.Load("testdata/yaml/aws-sigv4-profile-only.yml")
	assert.NoError(t, err)
	assert.NotNil(t, rn)

	r := rn.Spec.Requests[0]
	assert.NotNil(t, r.AWS)
	aws := r.GetAWSSign()
	assert.NotNil(t, aws)

	assert.Equal(t, "eu-west-1", aws.Region)
	assert.Equal(t, "default", aws.Profile)
}
