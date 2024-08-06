package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckBuild(t *testing.T) {
	Version = ""
	Date = ""
	Commit = ""

	CheckBuild()

	assert.Equal(t, notGiven, Version)
	assert.Equal(t, notGiven, Date)
	assert.Equal(t, notGiven, Commit)
}
