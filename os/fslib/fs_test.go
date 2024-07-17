package fslib

import (
	"testing"
)

func TestHasReadMode(t *testing.T) {
	HasReadMode("./fs.go")
}

func TestIsValidPath(t *testing.T) {
	if !IsValidPath("./fs.go") {
		t.Fatal("fs.go is valid path")
	}
}
