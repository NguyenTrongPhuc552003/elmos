// Package ui provides console output helpers for elmos.
package ui

import (
	"os"
	"testing"
)

func TestNewPrinter(t *testing.T) {
	p := NewPrinter()
	if p == nil {
		t.Fatal("NewPrinter() returned nil")
	}
}

func TestPrinter_Success(t *testing.T) {
	p := NewPrinter()
	// Test that Success doesn't panic
	p.Success("Task completed: %s", "build")
	p.Success("Simple message")
}

func TestPrinter_Error(t *testing.T) {
	p := NewPrinter()
	// Test that Error doesn't panic
	p.Error("Build failed: %s", "compile error")
	p.Error("Simple error")
}

func TestPrinter_Warn(t *testing.T) {
	p := NewPrinter()
	// Test that Warn doesn't panic
	p.Warn("Deprecated: %s", "old_function")
	p.Warn("Simple warning")
}

func TestPrinter_Info(t *testing.T) {
	p := NewPrinter()
	// Test that Info doesn't panic
	p.Info("Using arch: %s", "arm64")
	p.Info("Simple info")
}

func TestPrinter_Step(t *testing.T) {
	p := NewPrinter()
	// Test that Step doesn't panic
	p.Step("Building kernel...")
	p.Step("Compiling %s", "module.ko")
}

func TestPrinter_Print(t *testing.T) {
	p := NewPrinter()
	// Test that Print doesn't panic
	p.Print("Plain message")
	p.Print("Formatted: %d items", 5)
}

func TestPrinter_Writer(t *testing.T) {
	p := NewPrinter()
	w := p.Writer()
	if w != os.Stdout {
		t.Errorf("Printer.Writer() = %v, want os.Stdout", w)
	}
}

func TestPrintSuccess(t *testing.T) {
	// Test global function doesn't panic
	PrintSuccess("Global success: %s", "done")
	PrintSuccess("Simple message")
}

func TestPrintError(t *testing.T) {
	// Test global function doesn't panic
	PrintError("Global error: %s", "failed")
	PrintError("Simple error")
}

func TestPrintWarn(t *testing.T) {
	// Test global function doesn't panic
	PrintWarn("Global warning: %s", "caution")
	PrintWarn("Simple warning")
}

func TestPrintInfo(t *testing.T) {
	// Test global function doesn't panic
	PrintInfo("Global info: %s", "note")
	PrintInfo("Simple info")
}

func TestPrintStep(t *testing.T) {
	// Test global function doesn't panic
	PrintStep("Global step: %s", "building")
	PrintStep("Simple step")
}
