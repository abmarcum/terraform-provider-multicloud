package main

import (
	"testing"
)

func TestMainVersion(t *testing.T) {
	if version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", version)
	}
}
