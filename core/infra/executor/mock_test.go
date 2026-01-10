// Package executor provides abstractions for executing shell commands.
package executor

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestNewMockExecutor(t *testing.T) {
	got := NewMockExecutor()

	if got == nil {
		t.Fatal("NewMockExecutor() returned nil")
	}
	if got.Calls == nil {
		t.Error("NewMockExecutor().Calls is nil")
	}
	if len(got.Calls) != 0 {
		t.Errorf("NewMockExecutor().Calls length = %v, want 0", len(got.Calls))
	}
	if got.OutputResponses == nil {
		t.Error("NewMockExecutor().OutputResponses is nil")
	}
	if got.OutputErrors == nil {
		t.Error("NewMockExecutor().OutputErrors is nil")
	}
	if got.LookPathResponses == nil {
		t.Error("NewMockExecutor().LookPathResponses is nil")
	}
	if got.LookPathErrors == nil {
		t.Error("NewMockExecutor().LookPathErrors is nil")
	}
}

func TestMockExecutor_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			m: &MockExecutor{
				Calls: []CommandCall{},
			},
			args: args{
				ctx: context.Background(),
				cmd: "test",
			},
			wantErr: false,
		},
		{
			name: "Error",
			m: &MockExecutor{
				RunError: errors.New("mock error"),
			},
			args: args{
				ctx: context.Background(),
				cmd: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Run(tt.args.ctx, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("MockExecutor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockExecutor_RunWithEnv(t *testing.T) {
	type args struct {
		ctx  context.Context
		env  []string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			m:    NewMockExecutor(),
			args: args{
				ctx: context.Background(),
				env: []string{"KEY=VALUE"},
				cmd: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RunWithEnv(tt.args.ctx, tt.args.env, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("MockExecutor.RunWithEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockExecutor_RunInDir(t *testing.T) {
	type args struct {
		ctx  context.Context
		dir  string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			m:    NewMockExecutor(),
			args: args{
				ctx: context.Background(),
				dir: "/tmp",
				cmd: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RunInDir(tt.args.ctx, tt.args.dir, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("MockExecutor.RunInDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockExecutor_RunWithEnvInDir(t *testing.T) {
	type args struct {
		ctx  context.Context
		env  []string
		dir  string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			m:    NewMockExecutor(),
			args: args{
				ctx: context.Background(),
				env: []string{"KEY=VALUE"},
				dir: "/tmp",
				cmd: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RunWithEnvInDir(tt.args.ctx, tt.args.env, tt.args.dir, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("MockExecutor.RunWithEnvInDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockExecutor_Output(t *testing.T) {
	type args struct {
		ctx  context.Context
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Success",
			m: &MockExecutor{
				OutputResponses: map[string][]byte{"test": []byte("output")},
			},
			args: args{
				ctx: context.Background(),
				cmd: "test",
			},
			want:    []byte("output"),
			wantErr: false,
		},
		{
			name: "Error",
			m: &MockExecutor{
				OutputErrors: map[string]error{"test": errors.New("mock error")},
			},
			args: args{
				ctx: context.Background(),
				cmd: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Output(tt.args.ctx, tt.args.cmd, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MockExecutor.Output() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockExecutor.Output() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockExecutor_OutputWithEnv(t *testing.T) {
	type args struct {
		ctx  context.Context
		env  []string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Success",
			m: &MockExecutor{
				OutputResponses: map[string][]byte{"test": []byte("output")},
			},
			args: args{
				ctx: context.Background(),
				env: []string{"KEY=VALUE"},
				cmd: "test",
			},
			want:    []byte("output"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.OutputWithEnv(tt.args.ctx, tt.args.env, tt.args.cmd, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MockExecutor.OutputWithEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockExecutor.OutputWithEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockExecutor_LookPath(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			m: &MockExecutor{
				LookPathResponses: map[string]string{"test": "/bin/test"},
			},
			args: args{
				cmd: "test",
			},
			want:    "/bin/test",
			wantErr: false,
		},
		{
			name: "Error",
			m: &MockExecutor{
				LookPathErrors: map[string]error{"test": errors.New("mock error")},
			},
			args: args{
				cmd: "test",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.LookPath(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MockExecutor.LookPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("MockExecutor.LookPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockExecutor_Exec(t *testing.T) {
	type args struct {
		cmd  string
		args []string
		env  []string
	}
	tests := []struct {
		name    string
		m       *MockExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			m:    NewMockExecutor(),
			args: args{
				cmd: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Exec(tt.args.cmd, tt.args.args, tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("MockExecutor.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockExecutor_Reset(t *testing.T) {
	tests := []struct {
		name string
		m    *MockExecutor
	}{
		{
			name: "Reset",
			m: &MockExecutor{
				Calls: []CommandCall{{Cmd: "test"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Reset()
		})
	}
}
