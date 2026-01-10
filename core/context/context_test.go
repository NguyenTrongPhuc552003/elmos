// Package context provides build context management for elmos.
package context

import (
	"reflect"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg  *config.Config
		exec executor.Executor
		fs   filesystem.FileSystem
	}
	tests := []struct {
		name string
		args args
		want *Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg, tt.args.exec, tt.args.fs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_IsMounted(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.IsMounted(); got != tt.want {
				t.Errorf("Context.IsMounted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_EnsureMounted(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ctx.EnsureMounted(); (err != nil) != tt.wantErr {
				t.Errorf("Context.EnsureMounted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContext_GetActualMountPoint(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *Context
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ctx.GetActualMountPoint()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Context.GetActualMountPoint() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Context.GetActualMountPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_KernelExists(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.KernelExists(); got != tt.want {
				t.Errorf("Context.KernelExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_HasConfig(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.HasConfig(); got != tt.want {
				t.Errorf("Context.HasConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetKernelImage(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetKernelImage(); got != tt.want {
				t.Errorf("Context.GetKernelImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetVmlinux(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetVmlinux(); got != tt.want {
				t.Errorf("Context.GetVmlinux() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_HasKernelImage(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.HasKernelImage(); got != tt.want {
				t.Errorf("Context.HasKernelImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetMakeEnv(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetMakeEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Context.GetMakeEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_buildHostCFlags(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.buildHostCFlags(); got != tt.want {
				t.Errorf("Context.buildHostCFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_GetDefaultTargets(t *testing.T) {
	tests := []struct {
		name string
		ctx  *Context
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetDefaultTargets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Context.GetDefaultTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}
