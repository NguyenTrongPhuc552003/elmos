package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildModule(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildModule(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildModuleBuildCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildModuleBuildCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildModuleBuildCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildModuleListCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildModuleListCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildModuleListCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildModuleNewCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildModuleNewCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildModuleNewCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildModuleCleanCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildModuleCleanCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildModuleCleanCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildModuleHeaderCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildModuleHeaderCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildModuleHeaderCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOptionalArg(t *testing.T) {
	type args struct {
		args []string
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
			if got := getOptionalArg(tt.args.args); got != tt.want {
				t.Errorf("getOptionalArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
