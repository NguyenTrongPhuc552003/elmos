package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildPatch(t *testing.T) {
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
			if got := BuildPatch(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildPatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
