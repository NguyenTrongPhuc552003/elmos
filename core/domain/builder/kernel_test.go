package builder

import (
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
)

func TestBuildOptionsDefaults(t *testing.T) {
	opts := BuildOptions{}

	if opts.Jobs != 0 {
		t.Errorf("Default Jobs = %d, want 0", opts.Jobs)
	}
	if opts.Targets != nil {
		t.Errorf("Default Targets should be nil")
	}
}

func TestNewKernelBuilder(t *testing.T) {
	cfg := &elconfig.Config{
		Build: elconfig.BuildConfig{Arch: "arm64"},
	}

	// NewKernelBuilder should not panic with nil executor/fs
	// when used with proper mocks
	kb := NewKernelBuilder(nil, nil, cfg, nil)

	if kb == nil {
		t.Error("NewKernelBuilder returned nil")
	}
	if kb.cfg != cfg {
		t.Error("KernelBuilder.cfg not set correctly")
	}
}
