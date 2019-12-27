package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReading(t *testing.T) {
	bsIO, err := ReadFileFromPath("io.go")
	assert.NoError(t, err)
	assert.True(t, bytes.HasPrefix(bsIO, []byte("package")), "ReadFile io.go has no prefix 'package'")
}
