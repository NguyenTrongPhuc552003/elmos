// Package version provides version information for the elmos CLI.
package version

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		want Info
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_String(t *testing.T) {
	tests := []struct {
		name string
		i    Info
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("Info.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_Short(t *testing.T) {
	tests := []struct {
		name string
		i    Info
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Short(); got != tt.want {
				t.Errorf("Info.Short() = %v, want %v", got, tt.want)
			}
		})
	}
}
