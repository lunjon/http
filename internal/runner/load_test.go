package runner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/runner"
)

func TestLoadGoodJSON(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
	}{
		{"minimal json spec", "testdata/json/minimal.json"},
		{"post", "testdata/json/post.json"},
		{"headers", "testdata/json/headers.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := runner.LoadSpec(tt.filepath)
			assert.NoError(t, err)
			assert.NotNil(t, spec)
			err = spec.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestLoadBadJSON(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		failLoad bool
	}{
		{"ID with whitespace", "testdata/json/id_whitespace.json", false},
		{"wrong method json spec", "testdata/json/wrong_method.json", false},
		{"headers list", "testdata/json/headers_list.json", true},
		{"invalid URL json spec", "testdata/json/invalid_url.json", false},
		{"missing ID", "testdata/json/missing_id.json", false},
		{"missing URL", "testdata/json/missing_url.json", false},
		{"post, missing body", "testdata/json/post_missing_body.json", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := runner.LoadSpec(tt.filepath)
			if tt.failLoad {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, spec)
			}

			err = spec.Validate()
			assert.Error(t, err)
		})
	}
}

func TestLoadGoodYAML(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
	}{
		{"minimal yaml spec", "testdata/yaml/minimal.yaml"},
		{"minimal yml spec", "testdata/yaml/minimal.yml"},
		{"headers 1", "testdata/yaml/headers.yaml"},
		{"headers 2", "testdata/yaml/headers.yml"},
		{"post", "testdata/yaml/post.yml"},
		{"env", "testdata/yaml/env.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := runner.LoadSpec(tt.filepath)
			assert.NoError(t, err)
			assert.NotNil(t, spec)
			err = spec.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestLoadBadYAML(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		failLoad bool
	}{
		{"empty", "testdata/yaml/empty.yaml", false},
		{"post, missing body", "testdata/yaml/post_missing_body.yaml", false},
		{"env, wrong format", "testdata/yaml/env-wrong-format.yaml", false},
		{"unique names", "testdata/yaml/unique_ids.yaml", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := runner.LoadSpec(tt.filepath)
			if tt.failLoad {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, spec)
			}

			err = spec.Validate()
			assert.Error(t, err)
		})
	}
}

func TestLoadDefaultHeaders(t *testing.T) {
	spec, err := runner.LoadSpec("testdata/json/default_headers.json")
	assert.NoError(t, err)
	assert.NotNil(t, spec.Headers)

	err = spec.Validate()
	assert.NoError(t, err)

	r := spec.Requests[0]

	assert.Contains(t, r.Headers, "name")
	assert.Contains(t, r.Headers, "token")

	assert.Equal(t, r.Headers["name"], "override")
	assert.Equal(t, r.Headers["token"], "secret")
}

func TestLoadAWSSigv4Defaults(t *testing.T) {
	spec, err := runner.LoadSpec("testdata/yaml/aws-sigv4-bool.yml")
	assert.NoError(t, err)
	assert.NotNil(t, spec)

	err = spec.Validate()
	assert.NoError(t, err)

	r := spec.Requests[0]
	assert.NotNil(t, r.AWS)
	aws := r.GetAWSSign()
	assert.NotNil(t, aws)
	assert.Equal(t, constants.DefaultAWSRegion, aws.Region)
	assert.Empty(t, aws.Profile)
}

func TestLoadAWSSigv4RegionOnly(t *testing.T) {
	spec, err := runner.LoadSpec("testdata/yaml/aws-sigv4-profile-only.yml")
	assert.NoError(t, err)
	assert.NotNil(t, spec)

	err = spec.Validate()
	assert.NoError(t, err)

	r := spec.Requests[0]
	assert.NotNil(t, r.AWS)
	aws := r.GetAWSSign()
	assert.NotNil(t, aws)

	assert.Equal(t, constants.DefaultAWSRegion, aws.Region)
	assert.Equal(t, "default", aws.Profile)
}
