# Code Patterns

Go idioms and patterns used throughout ELMOS.

---

## Dependency Injection

ELMOS uses constructor injection for all domain services. The `App` struct wires everything:

```go
// core/app/app.go
func New(exec executor.Executor, fs filesystem.FileSystem, cfg *config.Config) *App {
    ctx := elcontext.New(cfg, exec, fs)
    printer := ui.NewPrinter()
    tm := toolchain.NewManager(exec, fs, cfg, printer)

    return &App{
        Exec:          exec,
        FS:            fs,
        Config:        cfg,
        Context:       ctx,
        KernelBuilder: builder.NewKernelBuilder(exec, fs, cfg, ctx, tm),
        ModuleBuilder: builder.NewModuleBuilder(exec, fs, cfg, ctx, tm),
        QEMURunner:    emulator.NewQEMURunner(exec, fs, cfg, ctx),
        // ...
    }
}
```

**Benefits:**

- Testable - inject mocks
- Explicit dependencies
- No global state

---

## Interface Abstraction

Domain defines interfaces, infra implements:

```go
// Domain expects:
type Executor interface {
    Run(ctx context.Context, name string, args ...string) error
    RunWithEnv(ctx context.Context, env []string, name string, args ...string) error
    Output(ctx context.Context, name string, args ...string) ([]byte, error)
}

// Infra provides:
type ShellExecutor struct{}  // Real implementation
type MockExecutor struct{}   // Test mock
```

This allows unit testing domain logic without real shell commands.

---

## Error Handling

### Custom Error Types

```go
// core/context/errors.go
type contextError struct {
    code    string
    message string
    cause   error
}

func (e *contextError) Error() string { return e.message }
func (e *contextError) Unwrap() error { return e.cause }
```

### Error Constructors

```go
func ImageError(msg string, cause error) error {
    return &contextError{code: "IMAGE", message: msg, cause: cause}
}

func ConfigError(msg string, cause error) error {
    return &contextError{code: "CONFIG", message: msg, cause: cause}
}
```

### Usage

```go
if !ctx.IsMounted() {
    return ImageError("kernel volume not mounted", ErrNotMounted)
}
```

---

## Configuration Pattern

### Struct with Tags

```go
// core/config/types.go
type Config struct {
    Image  ImageConfig  `mapstructure:"image"`
    Build  BuildConfig  `mapstructure:"build"`
    QEMU   QEMUConfig   `mapstructure:"qemu"`
    Paths  PathsConfig  `mapstructure:"paths"`
}
```

### Computed Defaults

```go
// core/config/loader.go
func applyComputedDefaults(cfg *Config) {
    if cfg.Paths.ProjectRoot == "" {
        cfg.Paths.ProjectRoot, _ = os.Getwd()
    }
    // Derive other paths from ProjectRoot...
}
```

---

## Embedded Assets

Templates embedded at compile time:

```go
// assets/embed.go
//go:embed templates/*
var Templates embed.FS

func GetModuleTemplate() ([]byte, error) {
    return Templates.ReadFile("templates/module/module.c.tmpl")
}
```

Used for scaffolding new modules/apps.

---

## Command Registration

### Grouped Commands

```go
// core/app/commands/register.go
func Register(ctx *Context, root *cobra.Command) {
    // Core commands
    root.AddCommand(BuildInit(ctx))
    root.AddCommand(BuildDoctor(ctx))
    
    // Build commands
    root.AddCommand(BuildKernel(ctx))
    root.AddCommand(BuildModule(ctx))
    
    // Runtime commands
    root.AddCommand(BuildQEMU(ctx))
}
```

### Command Context

All commands share a context:

```go
type Context struct {
    Exec          executor.Executor
    FS            filesystem.FileSystem
    Config        *config.Config
    KernelBuilder *builder.KernelBuilder
    Printer       *ui.Printer
    // ...
}
```

---

## Printer Pattern

Styled output with verbosity control:

```go
// core/ui/printer.go
type Printer struct{}

func (p *Printer) Step(format string, args ...interface{})    // → prefix
func (p *Printer) Success(format string, args ...interface{}) // ✓ prefix
func (p *Printer) Error(format string, args ...interface{})   // ✗ prefix
func (p *Printer) Info(format string, args ...interface{})    // ℹ prefix
func (p *Printer) Warn(format string, args ...interface{})    // ⚠ prefix
```

---

## Testing Patterns

### Table-Driven Tests

```go
func TestValidateMachine(t *testing.T) {
    tests := []struct {
        name    string
        machine string
        want    bool
    }{
        {"valid", "virt", true},
        {"invalid", "nonexistent", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

### Mock Executor

```go
type MockExecutor struct {
    outputs map[string][]byte
}

func (m *MockExecutor) Output(ctx context.Context, name string, args ...string) ([]byte, error) {
    return m.outputs[name], nil
}
```

---

## Naming Conventions

| Type        | Convention | Example                           |
| ----------- | ---------- | --------------------------------- |
| Interface   | Noun       | `Executor`, `FileSystem`          |
| Struct      | Noun       | `KernelBuilder`, `QEMURunner`     |
| Constructor | `New*`     | `NewKernelBuilder()`              |
| Method      | Verb       | `Build()`, `Run()`, `Configure()` |
| Error       | `*Error`   | `ImageError()`, `ConfigError()`   |