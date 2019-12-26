package util

import "testing"

import "bytes"

func TestReading(t *testing.T) {
	bsIO, err := ReadFileFromPath("io.go")
	if err != nil {
		t.Error("ReadFile io.go failed", err)
	}
	if !bytes.HasPrefix(bsIO, []byte("package")) {
		t.Error("ReadFile io.go has no prefix 'package'")
	}
}
