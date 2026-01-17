# Code Patterns

Go idioms and patterns used in ELMOS.

## Error Handling

Custom errors with context:

```go
type Error struct {
    Code    string
    Message string
    Cause   error
}

func (e Error) Error() string { return e.Message }
func (e Error) Unwrap() error { return e.Cause }
```

Used in `core/context/errors.go`.

## Dependency Injection

Infra interfaces injected into domain:

```go
type KernelBuilder struct {
    exec executor.Executor
    fs   filesystem.FileSystem
    // ...
}

func NewKernelBuilder(exec executor.Executor, fs filesystem.FileSystem, ...) *KernelBuilder {
    return &KernelBuilder{exec: exec, fs: fs}
}
```

## Interfaces for Testability

Domain defines interfaces, infra implements:

```go
type Executor interface {
    Run(cmd string, args ...string) error
    // ...
}
```

Mocks in `*_test.go` files.

## Configuration

YAML-based config with defaults:

```go
type Config struct {
    Arch     string `yaml:"arch"`
    Toolchain string `yaml:"toolchain"`
    // ...
}
```

Loaded via `config.Loader`.

## Embedded Assets

Templates embedded with `go:embed`:

```go
//go:embed templates/*
var templates embed.FS
```

Used for code generation.

## Cobra Commands

Structured commands in `commands/`:

```go
func init() {
    kernelCmd := &cobra.Command{Use: "kernel"}
    kernelCmd.AddCommand(&cobra.Command{Use: "clone"})
    // ...
}
```

## Logging

Simple printer with verbosity:

```go
type Printer struct{}

func (p *Printer) Printf(format string, args ...interface{})
func (p *Printer) Verbosef(format string, args ...interface{})
```

## Testing

Table-driven tests, mocks for infra.

## Naming

- Interfaces: `Executor`, `FileSystem`
- Structs: `KernelBuilder`, `QEMURunner`
- Methods: `Build()`, `Run()`