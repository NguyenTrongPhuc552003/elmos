// Package filesystem provides abstractions for file system operations.
package filesystem

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewOSFileSystem(t *testing.T) {
	tests := []struct {
		name string
		want *OSFileSystem
	}{
		{
			name: "Success",
			want: &OSFileSystem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOSFileSystem(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOSFileSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOSFileSystem_Stat(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exist.txt")
	if err := os.WriteFile(existingFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		wantErr bool
	}{
		{
			name:    "File Exists",
			f:       NewOSFileSystem(),
			args:    args{name: existingFile},
			wantErr: false,
		},
		{
			name:    "File Not Exists",
			f:       NewOSFileSystem(),
			args:    args{name: filepath.Join(tmpDir, "nonexistent.txt")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.f.Stat(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.Stat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOSFileSystem_ReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "read.txt")
	content := []byte("hello world")
	if err := os.WriteFile(existingFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Success",
			f:       NewOSFileSystem(),
			args:    args{name: existingFile},
			want:    content,
			wantErr: false,
		},
		{
			name:    "Not Found",
			f:       NewOSFileSystem(),
			args:    args{name: filepath.Join(tmpDir, "404.txt")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.ReadFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OSFileSystem.ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOSFileSystem_WriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "write.txt")

	type args struct {
		name string
		data []byte
		perm os.FileMode
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			f:    NewOSFileSystem(),
			args: args{
				name: targetFile,
				data: []byte("test data"),
				perm: 0644,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.WriteFile(tt.args.name, tt.args.data, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				content, err := os.ReadFile(tt.args.name)
				if err != nil {
					t.Fatalf("Failed to read written file: %v", err)
				}
				if !reflect.DeepEqual(content, tt.args.data) {
					t.Errorf("File content mismatch: got %v, want %v", content, tt.args.data)
				}
			}
		})
	}
}

func TestOSFileSystem_MkdirAll(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "a", "b", "c")

	type args struct {
		path string
		perm os.FileMode
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			f:    NewOSFileSystem(),
			args: args{
				path: targetDir,
				perm: 0755,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.MkdirAll(tt.args.path, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.MkdirAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				info, err := os.Stat(tt.args.path)
				if err != nil {
					t.Errorf("Directory not created: %v", err)
				}
				if !info.IsDir() {
					t.Errorf("Path is not a directory: %v", tt.args.path)
				}
			}
		})
	}
}

func TestOSFileSystem_ReadDir(t *testing.T) {
	tmpDir := t.TempDir()
	files := []string{"f1", "f2"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte(""), 0644); err != nil {
			t.Fatal(err)
		}
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		f         *OSFileSystem
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Success",
			f:         NewOSFileSystem(),
			args:      args{name: tmpDir},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Not Found",
			f:         NewOSFileSystem(),
			args:      args{name: filepath.Join(tmpDir, "missing")},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.ReadDir(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("OSFileSystem.ReadDir() count = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

func TestOSFileSystem_Remove(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "toremove.txt")
	if err := os.WriteFile(targetFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			f:       NewOSFileSystem(),
			args:    args{name: targetFile},
			wantErr: false,
		},
		{
			name:    "Not Found",
			f:       NewOSFileSystem(),
			args:    args{name: filepath.Join(tmpDir, "missing.txt")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Remove(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOSFileSystem_RemoveAll(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "toremove")
	if err := os.Mkdir(targetDir, 0755); err != nil {
		t.Fatal(err)
	}

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			f:       NewOSFileSystem(),
			args:    args{path: targetDir},
			wantErr: false,
		},
		{
			name:    "Success Even If Missing", // removeAll does not return error if path missing
			f:       NewOSFileSystem(),
			args:    args{path: filepath.Join(tmpDir, "missing")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.RemoveAll(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.RemoveAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOSFileSystem_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	existFile := filepath.Join(tmpDir, "exists")
	os.WriteFile(existFile, []byte(""), 0644)

	type args struct {
		path string
	}
	tests := []struct {
		name string
		f    *OSFileSystem
		args args
		want bool
	}{
		{
			name: "Exists",
			f:    NewOSFileSystem(),
			args: args{path: existFile},
			want: true,
		},
		{
			name: "Not Exists",
			f:    NewOSFileSystem(),
			args: args{path: filepath.Join(tmpDir, "missing")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Exists(tt.args.path); got != tt.want {
				t.Errorf("OSFileSystem.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOSFileSystem_IsDir(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "file")
	os.WriteFile(file, []byte(""), 0644)

	type args struct {
		path string
	}
	tests := []struct {
		name string
		f    *OSFileSystem
		args args
		want bool
	}{
		{
			name: "Is Directory",
			f:    NewOSFileSystem(),
			args: args{path: tmpDir},
			want: true,
		},
		{
			name: "Is File",
			f:    NewOSFileSystem(),
			args: args{path: file},
			want: false,
		},
		{
			name: "Not Exists",
			f:    NewOSFileSystem(),
			args: args{path: filepath.Join(tmpDir, "missing")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsDir(tt.args.path); got != tt.want {
				t.Errorf("OSFileSystem.IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOSFileSystem_Getwd(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		name    string
		f       *OSFileSystem
		want    string
		wantErr bool
	}{
		{
			name:    "Success",
			f:       NewOSFileSystem(),
			want:    wd,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Getwd()
			if (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.Getwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OSFileSystem.Getwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOSFileSystem_Create(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "create.txt")

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			f:       NewOSFileSystem(),
			args:    args{name: targetFile},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Create(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("OSFileSystem.Create() returned nil file")
			}
			if got != nil {
				got.Close()
			}
		})
	}
}

func TestOSFileSystem_Open(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "open.txt")
	os.WriteFile(targetFile, []byte(""), 0644)

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		f       *OSFileSystem
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			f:       NewOSFileSystem(),
			args:    args{name: targetFile},
			wantErr: false,
		},
		{
			name:    "Not Found",
			f:       NewOSFileSystem(),
			args:    args{name: filepath.Join(tmpDir, "missing")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("OSFileSystem.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("OSFileSystem.Open() returned nil file")
			}
			if got != nil {
				got.Close()
			}
		})
	}
}
