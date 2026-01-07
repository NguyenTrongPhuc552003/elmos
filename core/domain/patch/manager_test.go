package patch

import (
	"testing"
)

func TestPatchInfoFields(t *testing.T) {
	info := PatchInfo{
		Name:    "0001-fix-macos.patch",
		Path:    "/patches/v6.18/0001-fix-macos.patch",
		Version: "v6.18",
	}

	if info.Name != "0001-fix-macos.patch" {
		t.Errorf("PatchInfo.Name = %q, want %q", info.Name, "0001-fix-macos.patch")
	}
	if info.Version != "v6.18" {
		t.Errorf("PatchInfo.Version = %q, want %q", info.Version, "v6.18")
	}
}

func TestNewManager(t *testing.T) {
	m := NewManager(nil, nil, nil)

	if m == nil {
		t.Error("NewManager returned nil")
	}
}
