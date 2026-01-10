package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildApps(t *testing.T) {
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
			if got := BuildApps(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildApps() = %v, want %v", got, tt.want)
			}
		})
	}
}
