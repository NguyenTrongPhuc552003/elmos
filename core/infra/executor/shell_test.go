package executor

import (
	"os"
	"testing"
)

func TestNewShellExecutor(t *testing.T) {
	exec := NewShellExecutor()

	if exec == nil {
		t.Fatal("NewShellExecutor returned nil")
	}

	if exec.Stdout != os.Stdout {
		t.Error("Stdout should default to os.Stdout")
	}
	if exec.Stderr != os.Stderr {
		t.Error("Stderr should default to os.Stderr")
	}
	if exec.Stdin != os.Stdin {
		t.Error("Stdin should default to os.Stdin")
	}
}

func TestShellExecutorLookPath(t *testing.T) {
	exec := NewShellExecutor()

	// "ls" should exist on macOS/Linux
	path, err := exec.LookPath("ls")
	if err != nil {
		t.Skipf("ls not found, skipping test: %v", err)
	}
	if path == "" {
		t.Error("LookPath returned empty path for ls")
	}
}

func TestShellExecutorLookPathNotFound(t *testing.T) {
	exec := NewShellExecutor()

	// This command should not exist
	_, err := exec.LookPath("definitely-not-a-command-12345")
	if err == nil {
		t.Error("LookPath should return error for non-existent command")
	}
}

func TestShellExecutorImplementsInterface(t *testing.T) {
	// This test verifies the interface compile-time check
	var _ Executor = (*ShellExecutor)(nil)
}
