package executor

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Recorder is an Executor that runs nothing: it records every command and
// answers from canned results. Tests use it to check which host commands a
// domain operation would run, in which order, without touching the host.
type Recorder struct {
	mu sync.Mutex
	// Calls holds each command as "cmd arg1 arg2 ...", in call order.
	Calls []string
	// Outputs maps a command prefix ("hdiutil info") to the stdout it returns.
	Outputs map[string]string
	// Fail maps a command prefix to the error it returns.
	Fail map[string]error
	// NotOnPath lists commands that LookPath reports as missing.
	NotOnPath map[string]bool
}

// NewRecorder returns an empty Recorder.
func NewRecorder() *Recorder {
	return &Recorder{Outputs: map[string]string{}, Fail: map[string]error{}}
}

func (r *Recorder) record(cmd string, args []string) (string, error) {
	line := strings.TrimSpace(cmd + " " + strings.Join(args, " "))
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls = append(r.Calls, line)
	for prefix, err := range r.Fail {
		if strings.HasPrefix(line, prefix) {
			return "", err
		}
	}
	for prefix, out := range r.Outputs {
		if strings.HasPrefix(line, prefix) {
			return out, nil
		}
	}
	return "", nil
}

func (r *Recorder) Run(_ context.Context, cmd string, args ...string) error {
	_, err := r.record(cmd, args)
	return err
}

func (r *Recorder) RunWithEnv(_ context.Context, _ []string, cmd string, args ...string) error {
	_, err := r.record(cmd, args)
	return err
}

func (r *Recorder) RunInDir(_ context.Context, _ string, cmd string, args ...string) error {
	_, err := r.record(cmd, args)
	return err
}

func (r *Recorder) RunWithEnvInDir(_ context.Context, _ []string, _ string, cmd string, args ...string) error {
	_, err := r.record(cmd, args)
	return err
}

func (r *Recorder) Output(_ context.Context, cmd string, args ...string) ([]byte, error) {
	out, err := r.record(cmd, args)
	return []byte(out), err
}

func (r *Recorder) OutputWithEnv(_ context.Context, _ []string, cmd string, args ...string) ([]byte, error) {
	out, err := r.record(cmd, args)
	return []byte(out), err
}

func (r *Recorder) RunWithEnvSilent(_ context.Context, _ []string, cmd string, args ...string) error {
	_, err := r.record(cmd, args)
	return err
}

func (r *Recorder) LookPath(cmd string) (string, error) {
	if r.NotOnPath[cmd] {
		return "", fmt.Errorf("%s: not found", cmd)
	}
	return "/usr/bin/" + cmd, nil
}

func (r *Recorder) Exec(cmd string, args []string, _ []string) error {
	_, err := r.record(cmd, args)
	return err
}

var _ Executor = (*Recorder)(nil)
