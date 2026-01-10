// Package app provides the CLI application layer for elmos.
package app

import (
	"reflect"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/spf13/cobra"
)

func TestNew(t *testing.T) {
	type args struct {
		exec executor.Executor
		fs   filesystem.FileSystem
		cfg  *config.Config
	}
	tests := []struct {
		name string
		args args
		want *App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.exec, tt.args.fs, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_BuildRootCommand(t *testing.T) {
	tests := []struct {
		name string
		a    *App
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.BuildRootCommand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.BuildRootCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
