# Core App API

The app layer provides the CLI interface and dependency wiring.

## App Struct

```go
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
```

Holds all application dependencies, wired in `New()`.

## Key Functions

### New

```go
func New(exec executor.Executor, fs filesystem.FileSystem, cfg *config.Config) *App
```

Creates and wires all dependencies.

### BuildRootCommand

```go
func (a *App) BuildRootCommand() *cobra.Command
```

Builds the root Cobra command with all subcommands registered.

## Commands Context

```go
type Context struct {
    Exec             executor.Executor
    FS               filesystem.FileSystem
    Config           *config.Config
    AppContext       *elcontext.Context
    KernelBuilder    *builder.KernelBuilder
    // ... other builders
    Verbose          *bool
    ConfigFile       *string
}
```

Passed to command registrations for access to dependencies.

## Usage

The app layer is the entry point, initialized in `main.go` and executed via Cobra.