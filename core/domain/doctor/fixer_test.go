// Package doctor provides dependency checking and environment validation for elmos.
package doctor

import (
	"reflect"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNewAutoFixer(t *testing.T) {
	type args struct {
		fs  filesystem.FileSystem
		cfg *elconfig.Config
	}
	tests := []struct {
		name string
		args args
		want *AutoFixer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAutoFixer(tt.args.fs, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAutoFixer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutoFixer_FixElfH(t *testing.T) {
	tests := []struct {
		name    string
		f       *AutoFixer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.FixElfH(); (err != nil) != tt.wantErr {
				t.Errorf("AutoFixer.FixElfH() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAutoFixer_CanFixElfH(t *testing.T) {
	tests := []struct {
		name string
		f    *AutoFixer
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.CanFixElfH(); got != tt.want {
				t.Errorf("AutoFixer.CanFixElfH() = %v, want %v", got, tt.want)
			}
		})
	}
}
