// Package commands provides individual CLI command builders for elmos.
// Each command group is in its own subsystem package for maintainability.
package commands

import (
	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/build"
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/env"
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/ops"
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/runtime"
	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands/ui"
)

// Register adds all subcommands to the root command.
func Register(ctx *Context, rootCmd *cobra.Command) {
	// UI commands
	rootCmd.AddCommand(ui.BuildVersion(ctx))
	rootCmd.AddCommand(ui.BuildTUI(ctx))

	// Environment setup commands
	rootCmd.AddCommand(env.BuildInit(ctx))
	rootCmd.AddCommand(env.BuildExit(ctx))
	rootCmd.AddCommand(env.BuildArch(ctx))
	rootCmd.AddCommand(env.BuildToolchains(ctx))
	rootCmd.AddCommand(env.BuildRootfs(ctx))

	// Build commands
	rootCmd.AddCommand(build.BuildKernel(ctx))
	rootCmd.AddCommand(build.BuildModule(ctx))
	rootCmd.AddCommand(build.BuildApps(ctx))

	// Runtime commands
	rootCmd.AddCommand(runtime.BuildQEMU(ctx))
	rootCmd.AddCommand(runtime.BuildGDB(ctx))
	rootCmd.AddCommand(runtime.BuildPatch(ctx))
	rootCmd.AddCommand(runtime.BuildServe(ctx)) // gRPC API server

	// Operations commands
	rootCmd.AddCommand(ops.BuildDoctor(ctx))
	rootCmd.AddCommand(ops.BuildStatus(ctx))
}
