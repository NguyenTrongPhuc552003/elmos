package ui

import (
	"testing"
)

func TestNewPrinter(t *testing.T) {
	p := NewPrinter()
	if p == nil {
		t.Error("NewPrinter returned nil")
	}
}

func TestDefaultPrinterExists(t *testing.T) {
	// Verify the global defaultPrinter is initialized
	if defaultPrinter == nil {
		t.Error("defaultPrinter should not be nil")
	}
}
