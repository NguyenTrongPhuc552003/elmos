// Package app provides the CLI application layer for elmos.
package app

import (
	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/app/commands"
	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/builder"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/doctor"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/emulator"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/patch"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/rootfs"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/ui"
	"github.com/NguyenTrongPhuc552003/elmos/pkg/version"
)

// App holds all the application dependencies.
type App struct {
	Exec             executor.Executor
	FS               filesystem.FileSystem
	Config           *config.Config
	Context          *elcontext.Context
	KernelBuilder    *builder.KernelBuilder
	ModuleBuilder    *builder.ModuleBuilder
	AppBuilder       *builder.AppBuilder
	QEMURunner       *emulator.QEMURunner
	HealthChecker    *doctor.HealthChecker
	AutoFixer        *doctor.AutoFixer
	RootfsCreator    *rootfs.Creator
	PatchManager     *patch.Manager
	ToolchainManager *toolchain.Manager
	Printer          *ui.Printer
	Verbose          bool
	ConfigFile       string
}

// New creates a new App with all dependencies wired up.
func New(exec executor.Executor, fs filesystem.FileSystem, cfg *config.Config) *App {
	ctx := elcontext.New(cfg, exec, fs)
	return &App{
		Exec:             exec,
		FS:               fs,
		Config:           cfg,
		Context:          ctx,
		KernelBuilder:    builder.NewKernelBuilder(exec, fs, cfg, ctx),
		ModuleBuilder:    builder.NewModuleBuilder(exec, fs, cfg, ctx),
		AppBuilder:       builder.NewAppBuilder(exec, fs, cfg),
		QEMURunner:       emulator.NewQEMURunner(exec, fs, cfg, ctx),
		HealthChecker:    doctor.NewHealthChecker(exec, fs, cfg),
		AutoFixer:        doctor.NewAutoFixer(fs, cfg),
		RootfsCreator:    rootfs.NewCreator(exec, fs, cfg),
		PatchManager:     patch.NewManager(exec, fs, cfg),
		ToolchainManager: toolchain.NewManager(exec, fs, cfg),
		Printer:          ui.NewPrinter(),
	}
}

// BuildRootCommand builds the root cobra command with all subcommands.
func (a *App) BuildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "elmos",
		Short: "Embedded Linux on MacOS - Native kernel build tools",
		Long: `ELMOS provides native Linux kernel build tools for macOS.

Common workflow:
  elmos init              # Initialize workspace
  elmos doctor            # Check dependencies
  elmos kernel config     # Configure kernel
  elmos build             # Build kernel
  elmos qemu run          # Test in QEMU
  elmos tui               # Launch interactive TUI`,
		Version: version.Get().String(),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "version" || cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "tui" || cmd.Name() == "init" {
				return nil
			}
			// Reload config if custom path is provided
			if a.ConfigFile != "" {
				newCfg, err := config.Load(a.ConfigFile)
				if err != nil {
					return err
				}
				// update the struct contents so pointers passed to builders remain valid
				*a.Config = *newCfg
			}
			a.Context.Verbose = a.Verbose
			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&a.Verbose, "verbose", "e", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&a.ConfigFile, "config", "c", "", "config file (default is elmos.yaml)")

	// Create command context and register all commands
	cmdCtx := &commands.Context{
		Exec:             a.Exec,
		FS:               a.FS,
		Config:           a.Config,
		AppContext:       a.Context,
		KernelBuilder:    a.KernelBuilder,
		ModuleBuilder:    a.ModuleBuilder,
		AppBuilder:       a.AppBuilder,
		QEMURunner:       a.QEMURunner,
		HealthChecker:    a.HealthChecker,
		AutoFixer:        a.AutoFixer,
		RootfsCreator:    a.RootfsCreator,
		PatchManager:     a.PatchManager,
		ToolchainManager: a.ToolchainManager,
		Printer:          a.Printer,
		Verbose:          &a.Verbose,
		ConfigFile:       &a.ConfigFile,
	}

	commands.Register(cmdCtx, rootCmd)

	// Apply custom styled help output
	ui.SetCustomUsageFunc(rootCmd)

	return rootCmd
}
