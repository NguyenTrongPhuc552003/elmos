package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildQEMU(t *testing.T) {
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
			if got := BuildQEMU(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildQEMU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildGDB(t *testing.T) {
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
			if got := BuildGDB(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildGDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
