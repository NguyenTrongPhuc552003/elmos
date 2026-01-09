package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewOSFileSystem(t *testing.T) {
	fs := NewOSFileSystem()
	if fs == nil {
		t.Fatal("NewOSFileSystem returned nil")
	}
}

func TestOSFileSystemExists(t *testing.T) {
	fs := NewOSFileSystem()

	// Current directory should exist
	exists := fs.Exists(".")
	if !exists {
		t.Error("Exists(.) should return true")
	}

	// Non-existent path
	exists = fs.Exists("/definitely/not/a/real/path/12345")
	if exists {
		t.Error("Exists should return false for non-existent path")
	}
}

func TestOSFileSystemIsDir(t *testing.T) {
	fs := NewOSFileSystem()

	// Current directory is a directory
	if !fs.IsDir(".") {
		t.Error("IsDir(.) should return true")
	}

	// Non-existent path
	if fs.IsDir("/definitely/not/a/real/path/12345") {
		t.Error("IsDir should return false for non-existent path")
	}
}

func TestOSFileSystemGetwd(t *testing.T) {
	fs := NewOSFileSystem()

	cwd, err := fs.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error: %v", err)
	}
	if cwd == "" {
		t.Error("Getwd() returned empty string")
	}
}

func TestOSFileSystemReadWriteFile(t *testing.T) {
	fs := NewOSFileSystem()

	// Create temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	content := []byte("hello world")

	// Write file
	err := fs.WriteFile(tmpFile, content, 0644)
	if err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	// Read file
	data, err := fs.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("ReadFile() = %q, want %q", string(data), "hello world")
	}
}

func TestOSFileSystemMkdirAll(t *testing.T) {
	fs := NewOSFileSystem()

	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "a", "b", "c")

	err := fs.MkdirAll(newDir, 0755)
	if err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}

	if !fs.IsDir(newDir) {
		t.Error("Directory was not created")
	}
}

func TestOSFileSystemStat(t *testing.T) {
	fs := NewOSFileSystem()

	info, err := fs.Stat(".")
	if err != nil {
		t.Fatalf("Stat(.) error: %v", err)
	}
	if !info.IsDir() {
		t.Error("Stat(.) should return directory info")
	}
}

func TestOSFileSystemImplementsInterface(t *testing.T) {
	var _ FileSystem = (*OSFileSystem)(nil)
}

func TestOSFileSystemCreateOpen(t *testing.T) {
	fs := NewOSFileSystem()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "create_test.txt")

	// Create file
	f, err := fs.Create(tmpFile)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	_ = f.Close()

	// Open file
	f, err = fs.Open(tmpFile)
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	_ = f.Close()
}

func TestOSFileSystemRemove(t *testing.T) {
	fs := NewOSFileSystem()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "remove_test.txt")

	// Create file
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Remove file
	err := fs.Remove(tmpFile)
	if err != nil {
		t.Fatalf("Remove() error: %v", err)
	}

	if fs.Exists(tmpFile) {
		t.Error("File should be removed")
	}
}

func TestOSFileSystemRemoveAll(t *testing.T) {
	fs := NewOSFileSystem()

	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err := fs.RemoveAll(subDir)
	if err != nil {
		t.Fatalf("RemoveAll() error: %v", err)
	}

	if fs.Exists(subDir) {
		t.Error("Directory should be removed")
	}
}
