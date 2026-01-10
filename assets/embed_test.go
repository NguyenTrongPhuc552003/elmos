// Package assets provides embedded template files for elmos.
package assets

import (
	"reflect"
	"testing"
)

func TestGetModuleTemplate(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetModuleTemplate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetModuleTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetModuleTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetModuleMakefile(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetModuleMakefile()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetModuleMakefile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetModuleMakefile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppTemplate(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppTemplate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetAppTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppMakefile(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppMakefile()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetAppMakefile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppMakefile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInitScript(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInitScript()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetInitScript() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInitScript() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGuestSync(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGuestSync()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetGuestSync() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGuestSync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigTemplate(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfigTemplate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetConfigTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
