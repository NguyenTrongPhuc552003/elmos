// Package ui provides console output helpers for elmos.
package ui

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestBanner(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Banner(); got != tt.want {
				t.Errorf("Banner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCustomUsageFunc(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCustomUsageFunc(tt.args.cmd)
		})
	}
}

func Test_customHelpFunc(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customHelpFunc(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_writeHeader(t *testing.T) {
	type args struct {
		out *strings.Builder
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeHeader(tt.args.out, tt.args.cmd)
		})
	}
}

func Test_writeUsage(t *testing.T) {
	type args struct {
		out *strings.Builder
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeUsage(tt.args.out, tt.args.cmd)
		})
	}
}

func Test_writeCommands(t *testing.T) {
	type args struct {
		out *strings.Builder
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeCommands(tt.args.out, tt.args.cmd)
		})
	}
}

func Test_writeGroupedCommands(t *testing.T) {
	type args struct {
		out  *strings.Builder
		cmds []*cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeGroupedCommands(tt.args.out, tt.args.cmds)
		})
	}
}

func Test_writeSimpleCommands(t *testing.T) {
	type args struct {
		out  *strings.Builder
		cmds []*cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeSimpleCommands(tt.args.out, tt.args.cmds)
		})
	}
}

func Test_writeCommand(t *testing.T) {
	type args struct {
		out *strings.Builder
		sub *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeCommand(tt.args.out, tt.args.sub)
		})
	}
}

func Test_writeFlags(t *testing.T) {
	type args struct {
		out *strings.Builder
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeFlags(tt.args.out, tt.args.cmd)
		})
	}
}

func Test_writeExamples(t *testing.T) {
	type args struct {
		out *strings.Builder
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeExamples(tt.args.out, tt.args.cmd)
		})
	}
}

func Test_writeFooter(t *testing.T) {
	type args struct {
		out *strings.Builder
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeFooter(tt.args.out, tt.args.cmd)
		})
	}
}

func Test_customUsageFunc(t *testing.T) {
	type args struct {
		cmd *cobra.Command
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
			if err := customUsageFunc(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("customUsageFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupCommands(t *testing.T) {
	type args struct {
		cmds []*cobra.Command
	}
	tests := []struct {
		name string
		args args
		want []commandGroup
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := groupCommands(tt.args.cmds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildGroupSlice(t *testing.T) {
	type args struct {
		grouped map[string][]*cobra.Command
	}
	tests := []struct {
		name string
		args args
		want []commandGroup
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildGroupSlice(tt.args.grouped); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildGroupSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
