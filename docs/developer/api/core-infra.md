# Core Infra API

Infrastructure interfaces and implementations.

## Executor Interface

```go
type Executor interface {
    Run(cmd string, args ...string) error
    RunWithOutput(cmd string, args ...string) (string, error)
    RunInDir(dir, cmd string, args ...string) error
    RunAsync(cmd string, args ...string) (chan error, error)
}
```

Abstracts command execution.

### Implementations

- `ShellExecutor`: Real shell execution
- `MockExecutor`: For testing

## FileSystem Interface

```go
type FileSystem interface {
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, data []byte) error
    Exists(path string) bool
    MkdirAll(path string) error
    RemoveAll(path string) error
    // ... more methods
}
```

Abstracts file operations.

### Implementations

- `OSFileSystem`: Real OS operations
- `MockFileSystem`: In-memory for tests

## Homebrew Package

### Resolver

```go
type Resolver struct{}

func NewResolver() *Resolver
func (r *Resolver) IsInstalled(formula string) bool
func (r *Resolver) Install(formula string) error
```

Manages Homebrew dependencies.

## Usage

Infra provides pluggable implementations, allowing domain logic to be tested with mocks.