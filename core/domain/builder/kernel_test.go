// Package builder provides kernel and module build orchestration for elmos.
package builder

import (
	"context"
	"reflect"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNewKernelBuilder(t *testing.T) {
	type args struct {
		exec executor.Executor
		fs   filesystem.FileSystem
		cfg  *elconfig.Config
		ctx  *elcontext.Context
		tm   *toolchain.Manager
	}
	tests := []struct {
		name string
		args args
		want *KernelBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKernelBuilder(tt.args.exec, tt.args.fs, tt.args.cfg, tt.args.ctx, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKernelBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKernelBuilder_Build(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts BuildOptions
	}
	tests := []struct {
		name    string
		b       *KernelBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Build(tt.args.ctx, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("KernelBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKernelBuilder_Configure(t *testing.T) {
	type args struct {
		ctx        context.Context
		configType string
	}
	tests := []struct {
		name    string
		b       *KernelBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Configure(tt.args.ctx, tt.args.configType); (err != nil) != tt.wantErr {
				t.Errorf("KernelBuilder.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKernelBuilder_forceGraphicsConfig(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		b       *KernelBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.forceGraphicsConfig(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("KernelBuilder.forceGraphicsConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKernelBuilder_EnableKVMConfig(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		b       *KernelBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.EnableKVMConfig(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("KernelBuilder.EnableKVMConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKernelBuilder_Clean(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		b       *KernelBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Clean(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("KernelBuilder.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKernelBuilder_GetDefaultTargets(t *testing.T) {
	tests := []struct {
		name string
		b    *KernelBuilder
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetDefaultTargets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KernelBuilder.GetDefaultTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKernelBuilder_HasConfig(t *testing.T) {
	tests := []struct {
		name string
		b    *KernelBuilder
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.HasConfig(); got != tt.want {
				t.Errorf("KernelBuilder.HasConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKernelBuilder_HasKernelImage(t *testing.T) {
	tests := []struct {
		name string
		b    *KernelBuilder
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.HasKernelImage(); got != tt.want {
				t.Errorf("KernelBuilder.HasKernelImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
