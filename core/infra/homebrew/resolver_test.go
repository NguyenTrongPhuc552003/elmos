// Package homebrew provides utilities for resolving Homebrew package paths.
package homebrew

import (
	"errors"
	"reflect"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
)

func TestNewResolver(t *testing.T) {
	exec := executor.NewMockExecutor()
	r := NewResolver(exec)
	if r == nil {
		t.Fatal("NewResolver() returned nil")
	}
	if r.exec != exec {
		t.Error("NewResolver() did not set executor")
	}
}

func TestResolver_GetPrefix(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("/opt/homebrew/opt/llvm\n")

	r := NewResolver(exec)
	got := r.GetPrefix("llvm")
	want := "/opt/homebrew/opt/llvm"
	if got != want {
		t.Errorf("Resolver.GetPrefix() = %v, want %v", got, want)
	}

	// Test caching - should use cached value
	got2 := r.GetPrefix("llvm")
	if got2 != got {
		t.Error("GetPrefix() should cache results")
	}
}

func TestResolver_GetBin(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("/opt/homebrew/opt/llvm\n")

	r := NewResolver(exec)
	got := r.GetBin("llvm")
	want := "/opt/homebrew/opt/llvm/bin"
	if got != want {
		t.Errorf("Resolver.GetBin() = %v, want %v", got, want)
	}

	// Empty prefix returns empty bin
	exec2 := executor.NewMockExecutor()
	exec2.OutputErrors["brew"] = errors.New("not installed")
	r2 := NewResolver(exec2)
	got2 := r2.GetBin("nonexistent")
	if got2 != "" {
		t.Errorf("Resolver.GetBin() should return empty for error, got %v", got2)
	}
}

func TestResolver_GetSbin(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("/opt/homebrew/opt/e2fsprogs\n")

	r := NewResolver(exec)
	got := r.GetSbin("e2fsprogs")
	want := "/opt/homebrew/opt/e2fsprogs/sbin"
	if got != want {
		t.Errorf("Resolver.GetSbin() = %v, want %v", got, want)
	}
}

func TestResolver_GetInclude(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("/opt/homebrew/opt/libelf\n")

	r := NewResolver(exec)
	got := r.GetInclude("libelf")
	want := "/opt/homebrew/opt/libelf/include"
	if got != want {
		t.Errorf("Resolver.GetInclude() = %v, want %v", got, want)
	}
}

func TestResolver_GetLib(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("/opt/homebrew/opt/zlib\n")

	r := NewResolver(exec)
	got := r.GetLib("zlib")
	want := "/opt/homebrew/opt/zlib/lib"
	if got != want {
		t.Errorf("Resolver.GetLib() = %v, want %v", got, want)
	}
}

func TestResolver_GetLibexecBin(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("/opt/homebrew/opt/gnu-sed\n")

	r := NewResolver(exec)
	got := r.GetLibexecBin("gnu-sed")
	want := "/opt/homebrew/opt/gnu-sed/libexec/gnubin"
	if got != want {
		t.Errorf("Resolver.GetLibexecBin() = %v, want %v", got, want)
	}
}

func TestResolver_ListInstalled(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("llvm\nqemu\ngnu-sed\n")

	r := NewResolver(exec)
	got, err := r.ListInstalled()
	if err != nil {
		t.Fatalf("Resolver.ListInstalled() error = %v", err)
	}
	want := []string{"llvm", "qemu", "gnu-sed"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Resolver.ListInstalled() = %v, want %v", got, want)
	}

	// Test error case
	exec2 := executor.NewMockExecutor()
	exec2.OutputErrors["brew"] = errors.New("brew error")
	r2 := NewResolver(exec2)
	_, err = r2.ListInstalled()
	if err == nil {
		t.Error("Resolver.ListInstalled() should return error on brew failure")
	}
}

func TestResolver_ListTaps(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("homebrew/core\nhomebrew/cask\n")

	r := NewResolver(exec)
	got, err := r.ListTaps()
	if err != nil {
		t.Fatalf("Resolver.ListTaps() error = %v", err)
	}
	want := []string{"homebrew/core", "homebrew/cask"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Resolver.ListTaps() = %v, want %v", got, want)
	}
}

func TestResolver_IsInstalled(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("llvm\nqemu\n")

	r := NewResolver(exec)

	if !r.IsInstalled("llvm") {
		t.Error("IsInstalled() should return true for installed package")
	}
	if r.IsInstalled("nonexistent") {
		t.Error("IsInstalled() should return false for non-installed package")
	}
}

func TestResolver_IsTapped(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("homebrew/core\nhomebrew/cask\n")

	r := NewResolver(exec)

	if !r.IsTapped("homebrew/core") {
		t.Error("IsTapped() should return true for tapped repo")
	}
	if r.IsTapped("custom/tap") {
		t.Error("IsTapped() should return false for non-tapped repo")
	}
}

func TestResolver_ClearCache(t *testing.T) {
	exec := executor.NewMockExecutor()
	exec.OutputResponses["brew"] = []byte("/opt/homebrew/opt/llvm\n")

	r := NewResolver(exec)
	// Populate cache
	r.GetPrefix("llvm")

	// Clear cache
	r.ClearCache()

	// Cache should be empty (can't easily test, but shouldn't panic)
}
