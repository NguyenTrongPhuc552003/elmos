package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildInit(t *testing.T) {
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
			if got := BuildInit(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildExit(t *testing.T) {
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
			if got := BuildExit(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildExit() = %v, want %v", got, tt.want)
			}
		})
	}
}
