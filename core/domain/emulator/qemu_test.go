package emulator

import (
	"testing"
)

func TestContainsImpl(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"CONFIG_DEBUG_KERNEL=y", "DEBUG_KERNEL", true},
		{"hello world", "world", true},
		{"hello world", "foo", false},
		{"", "foo", false},
		{"abc", "abc", true},
		{"abc", "abcd", false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			if got := containsImpl(tt.s, tt.substr); got != tt.want {
				t.Errorf("containsImpl(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"CONFIG_DEBUG_KERNEL=y", "CONFIG_DEBUG_KERNEL=y", true},
		{"hello", "hello", true},
		{"hello world", "world", true},
		{"", "", true},
		{"a", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			if got := contains(tt.s, tt.substr); got != tt.want {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}
