package rootfs

import (
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
)

func TestGetDebianArch(t *testing.T) {
	tests := []struct {
		arch string
		want string
	}{
		{"arm64", "arm64"},
		{"arm", "armhf"},
		{"riscv", "riscv64"},
		{"unknown", "arm64"}, // Default fallback
	}

	for _, tt := range tests {
		t.Run(tt.arch, func(t *testing.T) {
			cfg := &config.Config{
				Build: config.BuildConfig{Arch: tt.arch},
			}
			c := &Creator{cfg: cfg}

			if got := c.getDebianArch(); got != tt.want {
				t.Errorf("getDebianArch() = %q, want %q", got, tt.want)
			}
		})
	}
}
