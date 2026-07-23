package main

import (
	"testing"
)

func TestMainVersion(t *testing.T) {
	if version != "0.1.0" {
		t.Errorf("expected version '0.1.0', got '%s'", version)
	}
}
