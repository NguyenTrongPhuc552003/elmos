// Package tui provides the interactive Text User Interface for elmos.
// This file contains the model initialization, menu structure, and entry point.
package tui

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewModel(t *testing.T) {
	tests := []struct {
		name string
		want Model
	}{
		{
			name: "Init Default Model",
			want: Model{
				width:      120,
				height:     30,
				leftWidth:  30,
				rightWidth: 90,
				// Other fields like menuStack, logLines are initialized empty
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewModel()
			opts := []cmp.Option{
				cmp.AllowUnexported(Model{}),
				cmpopts.EquateErrors(),
				// Ignore complex 3rd party structs and dynamic fields to avoid reflection issues
				cmpopts.IgnoreFields(Model{}, "spinner", "textInput", "viewport", "currentMenu", "execPath", "menuStack", "logLines"),
			}
			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Errorf("NewModel() mismatch (-want +got):\n%s", diff)
			}
			// Manual checks
			if len(got.currentMenu) == 0 {
				t.Error("NewModel() currentMenu is empty")
			}
			if got.execPath == "" {
				t.Error("NewModel() execPath is empty")
			}
			// Check viewport initialization
			if got.viewport.Width != 60 || got.viewport.Height != 20 {
				t.Errorf("NewModel() viewport dim = %dx%d, want 60x20", got.viewport.Width, got.viewport.Height)
			}
			// Check textinput initialization
			if got.textInput.Width != 40 {
				t.Errorf("NewModel() textInput width = %d, want 40", got.textInput.Width)
			}
		})
	}
}

func Test_buildMenuStructure(t *testing.T) {
	tests := []struct {
		name string
		want []MenuItem
	}{
		{
			name: "Menu Structure Integrity",
			want: nil, // We won't compare the huge struct, but just verify it returns something
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMenuStructure()
			if len(got) == 0 {
				t.Error("buildMenuStructure() returned empty menu")
			}
			// Verify key items exist
			foundKernel := false
			for _, item := range got {
				if item.Label == "Kernel" {
					foundKernel = true
					break
				}
			}
			if !foundKernel {
				t.Error("buildMenuStructure() missing Kernel menu")
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Run Interactive",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run() starts a full tea program which is hard to test in unit test without hijacking stdin/stdout.
			// Ideally we mock bubbletea program, but that's complex.
			// For now, we skip this or just ensure it doesn't panic if called (but it will block).
			// So we skip.
			t.Skip("Skipping interactive Run() test")
		})
	}
}

func Test_maxInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"First Larger", args{10, 5}, 10},
		{"Second Larger", args{3, 8}, 8},
		{"Equal", args{5, 5}, 5},
		{"Negative", args{-5, -2}, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("maxInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
