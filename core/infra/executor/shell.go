// Package executor provides abstractions for executing shell commands.
package executor

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"syscall"
)

// ShellExecutor implements Executor using os/exec.
type ShellExecutor struct {
	// Stdout is the writer for command stdout. Defaults to os.Stdout.
	Stdout *os.File
	// Stderr is the writer for command stderr. Defaults to os.Stderr.
	Stderr *os.File
	// Stdin is the reader for command stdin. Defaults to os.Stdin.
	Stdin *os.File
}

// NewShellExecutor creates a new ShellExecutor with default stdio.
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

// Run executes a command and waits for it to complete.
func (e *ShellExecutor) Run(ctx context.Context, cmd string, args ...string) error {
	return e.RunWithEnvInDir(ctx, nil, "", cmd, args...)
}

// RunWithEnv executes a command with custom environment variables.
func (e *ShellExecutor) RunWithEnv(ctx context.Context, env []string, cmd string, args ...string) error {
	return e.RunWithEnvInDir(ctx, env, "", cmd, args...)
}

// RunInDir executes a command in a specific working directory.
func (e *ShellExecutor) RunInDir(ctx context.Context, dir string, cmd string, args ...string) error {
	return e.RunWithEnvInDir(ctx, nil, dir, cmd, args...)
}

// RunWithEnvInDir executes a command with custom environment in a specific directory.
func (e *ShellExecutor) RunWithEnvInDir(ctx context.Context, env []string, dir string, cmd string, args ...string) error {
	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdout = e.Stdout
	c.Stderr = e.Stderr
	c.Stdin = e.Stdin

	if dir != "" {
		c.Dir = dir
	}

	if len(env) > 0 {
		// Merge custom env with current environment
		c.Env = append(os.Environ(), env...)
	}

	return c.Run()
}

// Output executes a command and returns its stdout.
func (e *ShellExecutor) Output(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	return e.OutputWithEnv(ctx, nil, cmd, args...)
}

// OutputWithEnv executes a command with custom environment and returns its stdout.
func (e *ShellExecutor) OutputWithEnv(ctx context.Context, env []string, cmd string, args ...string) ([]byte, error) {
	c := exec.CommandContext(ctx, cmd, args...)
	if len(env) > 0 {
		// Merge custom env with current environment
		c.Env = append(os.Environ(), env...)
	}
	return c.Output()
}

// RunWithEnvSilent executes a command with custom environment, suppressing stderr.
func (e *ShellExecutor) RunWithEnvSilent(ctx context.Context, env []string, cmd string, args ...string) error {
	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdout = e.Stdout
	c.Stderr = nil // Suppress stderr
	c.Stdin = e.Stdin

	if len(env) > 0 {
		c.Env = append(os.Environ(), env...)
	}

	return c.Run()
}

// LookPath searches for an executable in the system PATH.
func (e *ShellExecutor) LookPath(cmd string) (string, error) {
	return exec.LookPath(cmd)
}

// Exec replaces the current process with the specified command.
func (e *ShellExecutor) Exec(cmd string, args []string, env []string) error {
	return syscall.Exec(cmd, args, env)
}

// RunWithEnvStreaming executes a command and streams its output line-by-line.
func (e *ShellExecutor) RunWithEnvStreaming(ctx context.Context, env []string, cmd string, args ...string) (<-chan string, <-chan error) {
	linesCh := make(chan string, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(linesCh)
		defer close(errCh)

		c := exec.CommandContext(ctx, cmd, args...)

		if len(env) > 0 {
			c.Env = append(os.Environ(), env...)
		}

		// Create pipes for stdout and stderr
		stdout, err := c.StdoutPipe()
		if err != nil {
			errCh <- err
			return
		}

		stderr, err := c.StderrPipe()
		if err != nil {
			errCh <- err
			return
		}

		// Start the command
		if err := c.Start(); err != nil {
			errCh <- err
			return
		}

		// Create a combined reader for stdout and stderr
		// We'll use goroutines to read from both pipes
		done := make(chan struct{})

		// Read stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				select {
				case linesCh <- scanner.Text():
				case <-ctx.Done():
					return
				}
			}
		}()

		// Read stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				select {
				case linesCh <- scanner.Text():
				case <-ctx.Done():
					return
				}
			}
		}()

		// Wait for command to complete
		if err := c.Wait(); err != nil {
			errCh <- err
		}

		close(done)
	}()

	return linesCh, errCh
}

// Ensure ShellExecutor implements Executor.
var _ Executor = (*ShellExecutor)(nil)
