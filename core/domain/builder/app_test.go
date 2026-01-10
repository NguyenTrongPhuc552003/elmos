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

func TestNewAppBuilder(t *testing.T) {
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
		want *AppBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAppBuilder(tt.args.exec, tt.args.fs, tt.args.cfg, tt.args.ctx, tt.args.tm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppBuilder_Build(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		a       *AppBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Build(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("AppBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppBuilder_buildApp(t *testing.T) {
	type args struct {
		ctx      context.Context
		app      AppInfo
		compiler string
		env      []string
	}
	tests := []struct {
		name    string
		a       *AppBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.buildApp(tt.args.ctx, tt.args.app, tt.args.compiler, tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("AppBuilder.buildApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppBuilder_Clean(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		a       *AppBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Clean(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("AppBuilder.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppBuilder_GetApps(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		a       *AppBuilder
		args    args
		want    []AppInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.GetApps(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AppBuilder.GetApps() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppBuilder.GetApps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppBuilder_getSpecificApp(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		a       *AppBuilder
		args    args
		want    []AppInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.getSpecificApp(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AppBuilder.getSpecificApp() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppBuilder.getSpecificApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppBuilder_getAllApps(t *testing.T) {
	tests := []struct {
		name    string
		a       *AppBuilder
		want    []AppInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.getAllApps()
			if (err != nil) != tt.wantErr {
				t.Fatalf("AppBuilder.getAllApps() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppBuilder.getAllApps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppBuilder_getAppInfo(t *testing.T) {
	type args struct {
		name string
		path string
	}
	tests := []struct {
		name string
		a    *AppBuilder
		args args
		want AppInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.getAppInfo(tt.args.name, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppBuilder.getAppInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppBuilder_getCrossCompiler(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		a    *AppBuilder
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.getCrossCompiler(tt.args.prefix); got != tt.want {
				t.Errorf("AppBuilder.getCrossCompiler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppBuilder_CreateApp(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		a       *AppBuilder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.CreateApp(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("AppBuilder.CreateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_executeTemplate(t *testing.T) {
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
			got, err := executeTemplate(tt.args.name, tt.args.tmplContent, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("executeTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("executeTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
