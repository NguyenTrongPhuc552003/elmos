# Architecture

ELMOS follows a layered, domain-driven architecture in Go with strict separation of concerns.

## Visual Overview

For detailed diagrams, see the [Architecture Diagrams](diagrams.md) page:

- [Component Diagram](diagrams.md#component) - Package relationships
- [Sequence Diagram](diagrams.md#sequence) - Build workflow
- [Class Diagram](diagrams.md#class) - Domain types
- [State Diagram](diagrams.md#state) - Workspace lifecycle
- [Deployment Diagram](diagrams.md#deployment) - Runtime view

---

## Layer Overview

| Layer       | Package         | Purpose                            |
| ----------- | --------------- | ---------------------------------- |
| **App**     | `core/app/`     | CLI commands, dependency injection |
| **Domain**  | `core/domain/`  | Business logic, platform-agnostic  |
| **Infra**   | `core/infra/`   | System interactions, external deps |
| **Config**  | `core/config/`  | Configuration types and loading    |
| **Context** | `core/context/` | Build state and environment        |
| **UI**      | `core/ui/`      | TUI and output formatting          |

---

## App Layer (`core/app/`)

The application entry point and command registration.

### App Struct

```go
// core/app/app.go
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
}
```

**Key Functions:**

- `New()` - Creates App with all dependencies wired
- `BuildRootCommand()` - Returns the root Cobra command

### Commands Context

Commands receive a shared context:

```go
// core/app/commands/context.go
type Context struct {
    Exec             executor.Executor
    FS               filesystem.FileSystem
    Config           *config.Config
    AppContext       *elcontext.Context
    KernelBuilder    *builder.KernelBuilder
    // ... all domain services
    Printer          *ui.Printer
}
```

---

## Domain Layer (`core/domain/`)

Platform-independent business logic.

### Modules

| Module       | Purpose                    | Key Types                                      |
| ------------ | -------------------------- | ---------------------------------------------- |
| `builder/`   | Kernel, module, app builds | `KernelBuilder`, `ModuleBuilder`, `AppBuilder` |
| `doctor/`    | Environment health checks  | `HealthChecker`, `AutoFixer`                   |
| `emulator/`  | QEMU execution             | `QEMURunner`, `RunOptions`                     |
| `patch/`     | Kernel patch management    | `Manager`, `PatchInfo`                         |
| `rootfs/`    | Root filesystem creation   | `Creator`                                      |
| `toolchain/` | Cross-compiler management  | `Manager`                                      |

### Example: KernelBuilder

```go
// core/domain/builder/kernel.go
type KernelBuilder struct {
    exec executor.Executor
    fs   filesystem.FileSystem
    cfg  *config.Config
    ctx  *elcontext.Context
    tm   *toolchain.Manager
}

// Key methods
func (b *KernelBuilder) Build(ctx context.Context, opts BuildOptions) error
func (b *KernelBuilder) Configure(ctx context.Context, configType string) error
func (b *KernelBuilder) Clean(ctx context.Context) error
```

---

## Infra Layer (`core/infra/`)

External system interactions with mockable interfaces.

### Interfaces

| Package       | Interface    | Purpose                  |
| ------------- | ------------ | ------------------------ |
| `executor/`   | `Executor`   | Run shell commands       |
| `filesystem/` | `FileSystem` | File I/O operations      |
| `homebrew/`   | `Resolver`   | Homebrew path resolution |
| `printer/`    | `Printer`    | Formatted output         |

### Executor Interface

```go
// core/infra/executor/executor.go
type Executor interface {
    Run(ctx context.Context, name string, args ...string) error
    RunWithEnv(ctx context.Context, env []string, name string, args ...string) error
    Output(ctx context.Context, name string, args ...string) ([]byte, error)
    // ... more methods
}
```

---

## Context Layer (`core/context/`)

Build state and environment management.

### Context Struct

```go
// core/context/context.go
type Context struct {
    Config  *config.Config
    Exec    executor.Executor
    FS      filesystem.FileSystem
    Brew    *homebrew.Resolver
    Verbose bool
}
```

**Key Methods:**

| Method             | Purpose                              |
| ------------------ | ------------------------------------ |
| `IsMounted()`      | Check if workspace volume is mounted |
| `EnsureMounted()`  | Verify mount or return error         |
| `GetKernelImage()` | Path to built kernel image           |
| `GetMakeEnv()`     | Environment variables for make       |
| `HasConfig()`      | Check if `.config` exists            |

---

## Data Flow

```
User → CLI (Cobra) → App.BuildRootCommand()
                          ↓
                    Command Handler
                          ↓
              Domain Service (e.g., KernelBuilder)
                          ↓
              Infra Interface (e.g., Executor.Run())
                          ↓
                    External System (make, qemu, etc.)
```

1. User runs `elmos kernel build`
2. Cobra parses args, calls command handler
3. Handler uses `Context.KernelBuilder.Build()`
4. KernelBuilder calls `Executor.RunWithEnv()` with make args
5. Output streamed back to user

---

## Key Patterns

### Dependency Injection

All domain services receive dependencies via constructors:

```go
// app.go - wiring
func New(exec executor.Executor, fs filesystem.FileSystem, cfg *config.Config) *App {
    ctx := elcontext.New(cfg, exec, fs)
    tm := toolchain.NewManager(exec, fs, cfg, printer)
    
    return &App{
        KernelBuilder: builder.NewKernelBuilder(exec, fs, cfg, ctx, tm),
        // ...
    }
}
```

### Interface Abstraction

Domain defines interfaces, infra implements:

```go
// Domain uses:
type Executor interface { Run(...) error }

// Infra provides:
type ShellExecutor struct { /* implements Executor */ }
type MockExecutor struct  { /* for testing */ }
```

---

## Directory Structure

```
core/
├── app/
│   ├── app.go              # App struct, New(), BuildRootCommand()
│   ├── commands/           # CLI command handlers
│   │   ├── context.go      # Shared command context
│   │   ├── kernel.go       # elmos kernel *
│   │   ├── qemu.go         # elmos qemu *
│   │   └── ...
│   └── version/            # Version info
├── config/
│   ├── arch.go             # Architecture configs (arm64, riscv)
│   ├── defaults.go         # Default values
│   ├── loader.go           # YAML config loading
│   └── types.go            # Config struct definitions
├── context/
│   └── context.go          # Build context, mount checks
├── domain/
│   ├── builder/            # Kernel/module/app builders
│   ├── doctor/             # Health checks
│   ├── emulator/           # QEMU runner
│   ├── patch/              # Patch management
│   ├── rootfs/             # RootFS creation
│   └── toolchain/          # Toolchain management
├── infra/
│   ├── executor/           # Command execution
│   ├── filesystem/         # File operations
│   ├── homebrew/           # Homebrew paths
│   └── printer/            # Output formatting
└── ui/
    ├── printer.go          # Styled output
    └── tui/                 # Interactive TUI
```