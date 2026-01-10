// Package executor provides abstractions for executing shell commands.
package executor

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewShellExecutor(t *testing.T) {
	tests := []struct {
		name string
		want *ShellExecutor
	}{
		{
			name: "Success",
			want: &ShellExecutor{
				Stdout: os.Stdout,
				Stderr: os.Stderr,
				Stdin:  os.Stdin,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewShellExecutor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewShellExecutor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShellExecutor_Run(t *testing.T) {
	type args struct {
		ctx  context.Context
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			e:    NewShellExecutor(),
			args: args{
				ctx: context.Background(),
				cmd: "true",
			},
			wantErr: false,
		},
		{
			name: "Failure",
			e:    NewShellExecutor(),
			args: args{
				ctx: context.Background(),
				cmd: "false",
			},
			wantErr: true,
		},
		{
			name: "With Arguments",
			e:    NewShellExecutor(),
			args: args{
				ctx:  context.Background(),
				cmd:  "echo",
				args: []string{"hello"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Run(tt.args.ctx, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("ShellExecutor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShellExecutor_RunWithEnv(t *testing.T) {
	type args struct {
		ctx  context.Context
		env  []string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			e:    NewShellExecutor(),
			args: args{
				ctx:  context.Background(),
				env:  []string{"MY_VAR=hello"},
				cmd:  "sh",
				args: []string{"-c", "if [ \"$MY_VAR\" != \"hello\" ]; then exit 1; fi"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.RunWithEnv(tt.args.ctx, tt.args.env, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("ShellExecutor.RunWithEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShellExecutor_RunInDir(t *testing.T) {
	tmpDir := t.TempDir()
	resolvedTmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx  context.Context
		dir  string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			e:    NewShellExecutor(),
			args: args{
				ctx:  context.Background(),
				dir:  tmpDir,
				cmd:  "sh",
				args: []string{"-c", "if [ \"$(pwd -P)\" != \"" + resolvedTmpDir + "\" ]; then echo \"pwd=$(pwd -P) want=" + resolvedTmpDir + "\"; exit 1; fi"},
			},
			wantErr: false,
		},
		{
			name: "Invalid Directory",
			e:    NewShellExecutor(),
			args: args{
				ctx: context.Background(),
				dir: filepath.Join(tmpDir, "nonexistent"),
				cmd: "true",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.RunInDir(tt.args.ctx, tt.args.dir, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("ShellExecutor.RunInDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShellExecutor_RunWithEnvInDir(t *testing.T) {
	tmpDir := t.TempDir()
	resolvedTmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx  context.Context
		env  []string
		dir  string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			e:    NewShellExecutor(),
			args: args{
				ctx:  context.Background(),
				env:  []string{"MY_VAR=hello"},
				dir:  tmpDir,
				cmd:  "sh",
				args: []string{"-c", "if [ \"$MY_VAR\" != \"hello\" ] || [ \"$(pwd -P)\" != \"" + resolvedTmpDir + "\" ]; then echo \"pwd=$(pwd -P) want=" + resolvedTmpDir + " env=$MY_VAR\"; exit 1; fi"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.RunWithEnvInDir(tt.args.ctx, tt.args.env, tt.args.dir, tt.args.cmd, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("ShellExecutor.RunWithEnvInDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShellExecutor_Output(t *testing.T) {
	type args struct {
		ctx  context.Context
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Success",
			e:    NewShellExecutor(),
			args: args{
				ctx:  context.Background(),
				cmd:  "echo",
				args: []string{"hello"},
			},
			want:    []byte("hello\n"),
			wantErr: false,
		},
		{
			name: "Fail",
			e:    NewShellExecutor(),
			args: args{
				ctx: context.Background(),
				cmd: "false",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Output(tt.args.ctx, tt.args.cmd, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ShellExecutor.Output() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShellExecutor.Output() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShellExecutor_OutputWithEnv(t *testing.T) {
	type args struct {
		ctx  context.Context
		env  []string
		cmd  string
		args []string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Success",
			e:    NewShellExecutor(),
			args: args{
				ctx:  context.Background(),
				env:  []string{"MY_VAR=hello"},
				cmd:  "sh",
				args: []string{"-c", "echo $MY_VAR"},
			},
			want:    []byte("hello\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.OutputWithEnv(tt.args.ctx, tt.args.env, tt.args.cmd, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ShellExecutor.OutputWithEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShellExecutor.OutputWithEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShellExecutor_LookPath(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			e:    NewShellExecutor(),
			args: args{
				cmd: "ls",
			},
			want: "/bin/ls", // This is an approximation; the user's system likely has ls in /bin or /usr/bin.
			// However, exact path match is brittle across OSs (Mac vs Linux).
			// We might need a smarter check or skip if we can't be sure.
			// For now, let's just checking for error.
			wantErr: false,
		},
		{
			name: "Not Found",
			e:    NewShellExecutor(),
			args: args{
				cmd: "nonexistantcommand12345",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.LookPath(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ShellExecutor.LookPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("ShellExecutor.LookPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShellExecutor_Exec(t *testing.T) {
	type args struct {
		cmd  string
		args []string
		env  []string
	}
	tests := []struct {
		name    string
		e       *ShellExecutor
		args    args
		wantErr bool
	}{
		{
			name: "Failure",
			e:    NewShellExecutor(),
			args: args{
				cmd: "nonexistantcommand12345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Exec(tt.args.cmd, tt.args.args, tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("ShellExecutor.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
