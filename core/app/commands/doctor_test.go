package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildDoctor(t *testing.T) {
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
			if got := BuildDoctor(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildDoctor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSection(t *testing.T) {
	type args struct {
		name string
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
			if got := getSection(tt.args.name); got != tt.want {
				t.Errorf("getSection() = %v, want %v", got, tt.want)
			}
		})
	}
}
