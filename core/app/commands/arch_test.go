package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildArch(t *testing.T) {
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
			if got := BuildArch(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildArch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleArchNoArgs(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handleArchNoArgs(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("handleArchNoArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_showArchConfig(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := showArchConfig(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("showArchConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setArchTarget(t *testing.T) {
	type args struct {
		ctx    *Context
		cmd    *cobra.Command
		target string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setArchTarget(tt.args.ctx, tt.args.cmd, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("setArchTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resolveArchAndToolchain(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name                string
		args                args
		wantArchName        string
		wantToolchainTarget string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArchName, gotToolchainTarget := resolveArchAndToolchain(tt.args.target)
			if gotArchName != tt.wantArchName {
				t.Errorf("resolveArchAndToolchain() gotArchName = %v, want %v", gotArchName, tt.wantArchName)
			}
			if gotToolchainTarget != tt.wantToolchainTarget {
				t.Errorf("resolveArchAndToolchain() gotToolchainTarget = %v, want %v", gotToolchainTarget, tt.wantToolchainTarget)
			}
		})
	}
}

func Test_inferArchFromToolchain(t *testing.T) {
	type args struct {
		target string
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
			if got := inferArchFromToolchain(tt.args.target); got != tt.want {
				t.Errorf("inferArchFromToolchain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveArchConfig(t *testing.T) {
	type args struct {
		ctx      *Context
		archName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveArchConfig(tt.args.ctx, tt.args.archName); (err != nil) != tt.wantErr {
				t.Errorf("saveArchConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_selectToolchain(t *testing.T) {
	type args struct {
		ctx             *Context
		cmd             *cobra.Command
		toolchainTarget string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := selectToolchain(tt.args.ctx, tt.args.cmd, tt.args.toolchainTarget); (err != nil) != tt.wantErr {
				t.Errorf("selectToolchain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
