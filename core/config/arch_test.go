// Package config provides configuration management for elmos.
package config

import (
	"reflect"
	"testing"
)

func TestGetArchConfig(t *testing.T) {
	type args struct {
		arch string
	}
	tests := []struct {
		name string
		args args
		want *ArchConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetArchConfig(tt.args.arch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetArchConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSupportedArchitectures(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SupportedArchitectures(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupportedArchitectures() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidArch(t *testing.T) {
	type args struct {
		arch string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidArch(tt.args.arch); got != tt.want {
				t.Errorf("IsValidArch() = %v, want %v", got, tt.want)
			}
		})
	}
}
