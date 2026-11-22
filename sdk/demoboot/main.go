// Command demoboot is the Linux demo release's boot check
// (docs/developer/demo-release.md, workflow step 3): it prepares a release
// directory with the same code as `elmos demo linux` (demo.PrepareLinux),
// runs QEMU with demo.LinuxArgs, reads the console until every marker has
// appeared, and stops QEMU.
//
// Usage: go run ./sdk/demoboot DIR [SECONDS]
// Exit status 0 if every marker appeared within SECONDS (default 300).
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/NguyenTrongPhuc552003/elmos/core/domain/demo"
)

// Markers are the console lines the check waits for, in boot order.
var Markers = []string{
	"Linux version 6.18.0 (haven@ci)",
	"Run /init as init process",
	"[GUEST] kdal loaded OK",
	"KDAL resource demand",
	"System ready.",
}

// missing returns the markers not yet in seen.
func missing(seen map[string]bool) []string {
	var m []string
	for _, k := range Markers {
		if !seen[k] {
			m = append(m, k)
		}
	}
	return m
}

func run(dir string, limit time.Duration) error {
	if err := demo.PrepareLinux(dir); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, "qemu-system-aarch64", demo.LinuxArgs(dir)...)
	cmd.Stdin = nil
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout
	start := time.Now()
	if err := cmd.Start(); err != nil {
		return err
	}
	seen := map[string]bool{}
	sc := bufio.NewScanner(out)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		fmt.Println(line)
		for _, k := range Markers {
			if strings.Contains(line, k) {
				seen[k] = true
			}
		}
		if len(missing(seen)) == 0 {
			break
		}
	}
	_ = cmd.Process.Kill()
	_ = cmd.Wait()
	if m := missing(seen); len(m) > 0 {
		return fmt.Errorf("after %.0f s, markers missing: %s", time.Since(start).Seconds(), strings.Join(m, " | "))
	}
	fmt.Printf("demoboot: every marker in %.1f s\n", time.Since(start).Seconds())
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: demoboot DIR [SECONDS]")
		os.Exit(2)
	}
	limit := 300 * time.Second
	if len(os.Args) > 2 {
		s, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "demoboot: SECONDS is not a number")
			os.Exit(2)
		}
		limit = time.Duration(s) * time.Second
	}
	if err := run(os.Args[1], limit); err != nil {
		fmt.Fprintln(os.Stderr, "demoboot:", err)
		os.Exit(1)
	}
}
