package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildToolchains(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildToolchains(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildToolchains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildToolchainInstallCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildToolchainInstallCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildToolchainInstallCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildToolchainListCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildToolchainListCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildToolchainListCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildToolchainStatusCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildToolchainStatusCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildToolchainStatusCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_showToolchainStatus(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := showToolchainStatus(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("showToolchainStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buildToolchainBuildCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildToolchainBuildCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildToolchainBuildCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildToolchainMenuconfigCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildToolchainMenuconfigCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildToolchainMenuconfigCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildToolchainCleanCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildToolchainCleanCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildToolchainCleanCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildToolchainEnvCmd(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildToolchainEnvCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildToolchainEnvCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containsIgnoreCase(t *testing.T) {
	type args struct {
		s      string
		substr string
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
			if got := containsIgnoreCase(tt.args.s, tt.args.substr); got != tt.want {
				t.Errorf("containsIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
