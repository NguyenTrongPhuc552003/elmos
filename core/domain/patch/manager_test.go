// Package patch provides kernel patch management for elmos.
package patch

import (
	"context"
	"reflect"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNewManager(t *testing.T) {
	type args struct {
		exec executor.Executor
		fs   filesystem.FileSystem
		cfg  *elconfig.Config
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.exec, tt.args.fs, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Apply(t *testing.T) {
	type args struct {
		ctx       context.Context
		patchFile string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Apply(tt.args.ctx, tt.args.patchFile); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Reverse(t *testing.T) {
	type args struct {
		ctx       context.Context
		patchFile string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Reverse(tt.args.ctx, tt.args.patchFile); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Reverse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_List(t *testing.T) {
	tests := []struct {
		name    string
		m       *Manager
		want    []PatchInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.List()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Manager.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetPatchesForVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		want    []PatchInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetPatchesForVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Manager.GetPatchesForVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetPatchesForVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
