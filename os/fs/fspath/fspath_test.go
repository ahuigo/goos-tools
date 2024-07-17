package fspath

import (
	"testing"
)

func TestHasReadMode(t *testing.T) {
	StatHasReadMode("./fspath.go")
}

func TestIsValidPath(t *testing.T) {
	if !IsValidPath("./fspath.go") {
		t.Fatal("fs.go is valid path")
	}
}
