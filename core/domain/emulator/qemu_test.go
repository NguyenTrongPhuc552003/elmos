// Package emulator provides QEMU emulation orchestration for elmos.
package emulator

import (
	"context"
	"reflect"
	"testing"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestNewQEMURunner(t *testing.T) {
	type args struct {
		exec executor.Executor
		fs   filesystem.FileSystem
		cfg  *elconfig.Config
		ctx  *elcontext.Context
	}
	tests := []struct {
		name string
		args args
		want *QEMURunner
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQEMURunner(tt.args.exec, tt.args.fs, tt.args.cfg, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQEMURunner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQEMURunner_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts RunOptions
	}
	tests := []struct {
		name    string
		q       *QEMURunner
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.Run(tt.args.ctx, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("QEMURunner.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQEMURunner_Debug(t *testing.T) {
	type args struct {
		ctx       context.Context
		graphical bool
	}
	tests := []struct {
		name    string
		q       *QEMURunner
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.Debug(tt.args.ctx, tt.args.graphical); (err != nil) != tt.wantErr {
				t.Errorf("QEMURunner.Debug() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQEMURunner_ConnectGDB(t *testing.T) {
	tests := []struct {
		name    string
		q       *QEMURunner
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.ConnectGDB(); (err != nil) != tt.wantErr {
				t.Errorf("QEMURunner.ConnectGDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQEMURunner_buildArgs(t *testing.T) {
	type args struct {
		archCfg     *elconfig.ArchConfig
		kernelImage string
		opts        RunOptions
	}
	tests := []struct {
		name string
		q    *QEMURunner
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.buildArgs(tt.args.archCfg, tt.args.kernelImage, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QEMURunner.buildArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQEMURunner_executeQEMU(t *testing.T) {
	type args struct {
		ctx    context.Context
		binary string
		args   []string
	}
	tests := []struct {
		name    string
		q       *QEMURunner
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.executeQEMU(tt.args.ctx, tt.args.binary, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("QEMURunner.executeQEMU() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQEMURunner_prepareModulesSync(t *testing.T) {
	tests := []struct {
		name    string
		q       *QEMURunner
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.prepareModulesSync(); (err != nil) != tt.wantErr {
				t.Errorf("QEMURunner.prepareModulesSync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQEMURunner_CheckDebugConfig(t *testing.T) {
	tests := []struct {
		name    string
		q       *QEMURunner
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.CheckDebugConfig(); (err != nil) != tt.wantErr {
				t.Errorf("QEMURunner.CheckDebugConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
