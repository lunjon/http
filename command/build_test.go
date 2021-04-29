package command_test

import (
	"github.com/lunjon/http/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	http := command.Build("0.1.0")
	assert.NotNil(t, http)
}
