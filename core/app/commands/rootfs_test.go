package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildRootfs(t *testing.T) {
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
			if got := BuildRootfs(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildRootfs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatBytes(t *testing.T) {
	type args struct {
		b int64
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
			if got := formatBytes(tt.args.b); got != tt.want {
				t.Errorf("formatBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
