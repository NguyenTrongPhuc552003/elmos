package terminal

import (
	"bytes"
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// popMenuStack navigates back in the menu hierarchy.
func (m *Model) popMenuStack() {
	if len(m.menuStack) > 0 {
		m.currentMenu = m.menuStack[len(m.menuStack)-1]
		m.menuStack = m.menuStack[:len(m.menuStack)-1]
		m.cursor, m.parentTitle = 0, ""
	}
}

// commandFormatters maps action identifiers to command format strings.
var commandFormatters = map[string]string{
	"module:new":           "elmos module new %s",
	"module:build:one":     "elmos module build %s",
	"app:new":              "elmos app new %s",
	"app:build:one":        "elmos app build %s",
	"config:arch":          "elmos config set arch %s",
	"config:jobs":          "elmos config set jobs %s",
	"config:memory":        "elmos config set memory %s",
	"rootfs:create:custom": "elmos rootfs create -s %s",
	"toolchain:select":     "elmos toolchains %s",
	"kernel:switch":        "elmos kernel switch %s",
}

// getCommandWithInput returns the display command string for a given action and input.
func (m *Model) getCommandWithInput(action, value string) string {
	if format, ok := commandFormatters[action]; ok {
		return fmt.Sprintf(format, value)
	}
	return "elmos " + action
}

// runCommand executes a command asynchronously and returns the result.
func (m *Model) runCommand(action string, args []string) tea.Cmd {
	return func() tea.Msg {
		// args passed directly
		cmd := exec.Command(m.execPath, args...)
		var output bytes.Buffer
		cmd.Stdout, cmd.Stderr = &output, &output
		err := cmd.Run()
		return CommandDoneMsg{Action: action, Err: err, Output: output.String()}
	}
}

// actionArgsDispatch maps action identifiers to argument generators.
// Only dynamic actions need to be here now.
var actionArgsDispatch = map[string]func(string) []string{
	"arch:set": func(v string) []string { return []string{"arch", v} },
	"kernel:switch": func(v string) []string {
		if v == "" {
			return []string{"kernel", "switch"}
		}
		return []string{"kernel", "switch", v}
	},
	"kernel:config": func(v string) []string {
		if v == "" || v == "defconfig" {
			return []string{"kernel", "config"}
		}
		return []string{"kernel", "config", v}
	},
	"module:build": func(v string) []string {
		if v == "" {
			return []string{"module", "build"}
		}
		return []string{"module", "build", v}
	},
	"module:new": func(v string) []string { return []string{"module", "new", v} },
	"app:build": func(v string) []string {
		if v == "" {
			return []string{"app", "build"}
		}
		return []string{"app", "build", v}
	},
	"app:new":              func(v string) []string { return []string{"app", "new", v} },
	"rootfs:create:custom": func(v string) []string { return []string{"rootfs", "create", "-s", v} },
	"config:arch":          func(v string) []string { return []string{"config", "set", "arch", v} },
	"config:jobs":          func(v string) []string { return []string{"config", "set", "jobs", v} },
	"config:memory":        func(v string) []string { return []string{"config", "set", "memory", v} },
	"toolchain:select":     func(v string) []string { return []string{"toolchains", v} },
}

// actionToArgs converts an action identifier to CLI arguments using map dispatch.
func (m *Model) actionToArgs(action, inputValue string) []string {
	if fn, ok := actionArgsDispatch[action]; ok {
		return fn(inputValue)
	}
	return []string{}
}

// isInteractiveCommand checks if a command from input mode should be run interactively.
func (m *Model) isInteractiveCommand(action, value string) bool {
	if action == "kernel:config" {
		// menuconfig, nconfig, xconfig, gconfig need a TTY
		return value == "menuconfig" || value == "nconfig" || value == "xconfig" || value == "gconfig"
	}
	return false
}
