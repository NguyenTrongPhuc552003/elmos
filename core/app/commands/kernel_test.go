package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildKernel(t *testing.T) {
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
			if got := BuildKernel(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildKernel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelConfigCmd(t *testing.T) {
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
			if got := buildKernelConfigCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelConfigCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelCleanCmd(t *testing.T) {
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
			if got := buildKernelCleanCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelCleanCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelCloneCmd(t *testing.T) {
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
			if got := buildKernelCloneCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelCloneCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelStatusCmd(t *testing.T) {
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
			if got := buildKernelStatusCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelStatusCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelResetCmd(t *testing.T) {
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
			if got := buildKernelResetCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelResetCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelSwitchCmd(t *testing.T) {
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
			if got := buildKernelSwitchCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelSwitchCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelPullCmd(t *testing.T) {
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
			if got := buildKernelPullCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelPullCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildKernelBuildCmd(t *testing.T) {
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
			if got := buildKernelBuildCmd(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKernelBuildCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printKernelGitInfo(t *testing.T) {
	type args struct {
		ctx *Context
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printKernelGitInfo(tt.args.ctx, tt.args.cmd)
		})
	}
}

func Test_printKernelBuildStatus(t *testing.T) {
	type args struct {
		ctx *Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printKernelBuildStatus(tt.args.ctx)
		})
	}
}

func Test_listKernelRefs(t *testing.T) {
	type args struct {
		ctx *Context
		cmd *cobra.Command
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
			if err := listKernelRefs(tt.args.ctx, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("listKernelRefs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_switchKernelRef(t *testing.T) {
	type args struct {
		ctx *Context
		cmd *cobra.Command
		ref string
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
			if err := switchKernelRef(tt.args.ctx, tt.args.cmd, tt.args.ref); (err != nil) != tt.wantErr {
				t.Errorf("switchKernelRef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
