// Package rootfs provides rootfs creation functionality for elmos.
package rootfs

import (
	"context"
	"reflect"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNewCreator(t *testing.T) {
	type args struct {
		exec executor.Executor
		fs   filesystem.FileSystem
		cfg  *elconfig.Config
	}
	tests := []struct {
		name string
		args args
		want *Creator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCreator(tt.args.exec, tt.args.fs, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCreator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreator_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts CreateOptions
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Create(tt.args.ctx, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Creator.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_cleanRootfsDir(t *testing.T) {
	type args struct {
		ctx       context.Context
		rootfsDir string
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.cleanRootfsDir(tt.args.ctx, tt.args.rootfsDir); (err != nil) != tt.wantErr {
				t.Errorf("Creator.cleanRootfsDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_runDebootstrap(t *testing.T) {
	type args struct {
		ctx       context.Context
		rootfsDir string
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.runDebootstrap(tt.args.ctx, tt.args.rootfsDir); (err != nil) != tt.wantErr {
				t.Errorf("Creator.runDebootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_removeDiskImage(t *testing.T) {
	type args struct {
		ctx       context.Context
		diskImage string
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.removeDiskImage(tt.args.ctx, tt.args.diskImage); (err != nil) != tt.wantErr {
				t.Errorf("Creator.removeDiskImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_createDiskImage(t *testing.T) {
	type args struct {
		ctx       context.Context
		diskImage string
		rootfsDir string
		size      string
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.createDiskImage(tt.args.ctx, tt.args.diskImage, tt.args.rootfsDir, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("Creator.createDiskImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_fixAptLists(t *testing.T) {
	type args struct {
		rootfsDir string
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.fixAptLists(tt.args.rootfsDir); (err != nil) != tt.wantErr {
				t.Errorf("Creator.fixAptLists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_getDebianArch(t *testing.T) {
	tests := []struct {
		name string
		c    *Creator
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.getDebianArch(); got != tt.want {
				t.Errorf("Creator.getDebianArch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreator_createInitScript(t *testing.T) {
	type args struct {
		rootfsDir string
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.createInitScript(tt.args.rootfsDir); (err != nil) != tt.wantErr {
				t.Errorf("Creator.createInitScript() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_Status(t *testing.T) {
	tests := []struct {
		name    string
		c       *Creator
		want    *RootfsInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Status()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Creator.Status() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Creator.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreator_Clean(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *Creator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Clean(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Creator.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreator_Exists(t *testing.T) {
	tests := []struct {
		name string
		c    *Creator
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Exists(); got != tt.want {
				t.Errorf("Creator.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
