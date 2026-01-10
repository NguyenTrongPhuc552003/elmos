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

func TestNewModuleBuilder(t *testing.T) {
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
		want *ModuleBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewModuleBuilder(tt.args.exec, tt.args.fs, tt.args.cfg, tt.args.ctx, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModuleBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleBuilder_Build(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		m       *ModuleBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Build(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ModuleBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModuleBuilder_buildModule(t *testing.T) {
	type args struct {
		ctx context.Context
		mod ModuleInfo
	}
	tests := []struct {
		name    string
		m       *ModuleBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.buildModule(tt.args.ctx, tt.args.mod); (err != nil) != tt.wantErr {
				t.Errorf("ModuleBuilder.buildModule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModuleBuilder_Clean(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		m       *ModuleBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Clean(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ModuleBuilder.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModuleBuilder_GetModules(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *ModuleBuilder
		args    args
		want    []ModuleInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetModules(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ModuleBuilder.GetModules() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleBuilder.GetModules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleBuilder_getSpecificModule(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *ModuleBuilder
		args    args
		want    []ModuleInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.getSpecificModule(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ModuleBuilder.getSpecificModule() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleBuilder.getSpecificModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleBuilder_getAllModules(t *testing.T) {
	tests := []struct {
		name    string
		m       *ModuleBuilder
		want    []ModuleInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.getAllModules()
			if (err != nil) != tt.wantErr {
				t.Fatalf("ModuleBuilder.getAllModules() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleBuilder.getAllModules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleBuilder_getModuleInfo(t *testing.T) {
	type args struct {
		name string
		path string
	}
	tests := []struct {
		name string
		m    *ModuleBuilder
		args args
		want ModuleInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.getModuleInfo(tt.args.name, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleBuilder.getModuleInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractModuleDescription(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractModuleDescription(tt.args.content); got != tt.want {
				t.Errorf("extractModuleDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleBuilder_PrepareHeaders(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       *ModuleBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.PrepareHeaders(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ModuleBuilder.PrepareHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModuleBuilder_CreateModule(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		m       *ModuleBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.CreateModule(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ModuleBuilder.CreateModule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_executeModuleTemplate(t *testing.T) {
	type args struct {
		name        string
		tmplContent string
		data        interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executeModuleTemplate(tt.args.name, tt.args.tmplContent, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("executeModuleTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("executeModuleTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
