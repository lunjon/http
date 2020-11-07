package command_test

import (
	"github.com/lunjon/httpreq/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	httpreq := command.Build()
	assert.NotNil(t, httpreq)
}
