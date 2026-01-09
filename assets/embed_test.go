package assets

import (
	"testing"
)

func TestGetModuleTemplate(t *testing.T) {
	data, err := GetModuleTemplate()
	if err != nil {
		t.Fatalf("GetModuleTemplate() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("GetModuleTemplate() returned empty data")
	}
}

func TestGetModuleMakefile(t *testing.T) {
	data, err := GetModuleMakefile()
	if err != nil {
		t.Fatalf("GetModuleMakefile() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("GetModuleMakefile() returned empty data")
	}
}

func TestGetAppTemplate(t *testing.T) {
	data, err := GetAppTemplate()
	if err != nil {
		t.Fatalf("GetAppTemplate() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("GetAppTemplate() returned empty data")
	}
}

func TestGetAppMakefile(t *testing.T) {
	data, err := GetAppMakefile()
	if err != nil {
		t.Fatalf("GetAppMakefile() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("GetAppMakefile() returned empty data")
	}
}

func TestGetInitScript(t *testing.T) {
	data, err := GetInitScript()
	if err != nil {
		t.Fatalf("GetInitScript() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("GetInitScript() returned empty data")
	}
}

func TestGetGuestSync(t *testing.T) {
	data, err := GetGuestSync()
	if err != nil {
		t.Fatalf("GetGuestSync() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("GetGuestSync() returned empty data")
	}
}
