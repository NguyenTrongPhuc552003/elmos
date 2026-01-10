// Package homebrew provides utilities for resolving Homebrew package paths.
package homebrew

import (
	"reflect"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
)

func TestNewResolver(t *testing.T) {
	type args struct {
		exec executor.Executor
	}
	tests := []struct {
		name string
		args args
		want *Resolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResolver(tt.args.exec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResolver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_GetPrefix(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetPrefix(tt.args.pkg); got != tt.want {
				t.Errorf("Resolver.GetPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_GetBin(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetBin(tt.args.pkg); got != tt.want {
				t.Errorf("Resolver.GetBin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_GetSbin(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetSbin(tt.args.pkg); got != tt.want {
				t.Errorf("Resolver.GetSbin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_GetInclude(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetInclude(tt.args.pkg); got != tt.want {
				t.Errorf("Resolver.GetInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_GetLib(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetLib(tt.args.pkg); got != tt.want {
				t.Errorf("Resolver.GetLib() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_GetLibexecBin(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetLibexecBin(tt.args.pkg); got != tt.want {
				t.Errorf("Resolver.GetLibexecBin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_ListInstalled(t *testing.T) {
	tests := []struct {
		name    string
		r       *Resolver
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ListInstalled()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Resolver.ListInstalled() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver.ListInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_ListTaps(t *testing.T) {
	tests := []struct {
		name    string
		r       *Resolver
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ListTaps()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Resolver.ListTaps() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver.ListTaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_IsInstalled(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsInstalled(tt.args.pkg); got != tt.want {
				t.Errorf("Resolver.IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_IsTapped(t *testing.T) {
	type args struct {
		tap string
	}
	tests := []struct {
		name string
		r    *Resolver
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsTapped(tt.args.tap); got != tt.want {
				t.Errorf("Resolver.IsTapped() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_ClearCache(t *testing.T) {
	tests := []struct {
		name string
		r    *Resolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.ClearCache()
		})
	}
}
