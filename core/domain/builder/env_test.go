package builder

import (
	"reflect"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func Test_getToolchainEnv(t *testing.T) {
	type args struct {
		ctx  *elcontext.Context
		cfg  *elconfig.Config
		tm   *toolchain.Manager
		fs   filesystem.FileSystem
		arch string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getToolchainEnv(tt.args.ctx, tt.args.cfg, tt.args.tm, tt.args.fs, tt.args.arch)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getToolchainEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getToolchainEnv() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getToolchainEnv() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
