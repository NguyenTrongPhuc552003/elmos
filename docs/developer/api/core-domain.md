# Core Domain API

Business logic modules in the domain layer.

## Builder Package

### KernelBuilder

```go
type KernelBuilder struct {
    // Fields for kernel operations
}

func NewKernelBuilder(...) *KernelBuilder
func (kb *KernelBuilder) Clone() error
func (kb *KernelBuilder) Configure(configType string) error
func (kb *KernelBuilder) Build() error
func (kb *KernelBuilder) Clean() error
```

Handles kernel cloning, config, building.

### ModuleBuilder / AppBuilder

Similar structure for modules and apps.

- `Create(name string)` - Generate templates
- `Build(dir string)` - Cross-compile

## Doctor Package

### HealthChecker

```go
type HealthChecker struct{}

func NewHealthChecker(...) *HealthChecker
func (hc *HealthChecker) Check() ([]CheckResult, error)
```

Runs dependency checks.

### AutoFixer

```go
type AutoFixer struct{}

func NewAutoFixer(...) *AutoFixer
func (af *AutoFixer) Fix(results []CheckResult) error
```

Applies automatic fixes.

## Emulator Package

### QEMURunner

```go
type QEMURunner struct{}

func NewQEMURunner(...) *QEMURunner
func (qr *QEMURunner) Run(options []string) error
func (qr *QEMURunner) Debug(options []string) error
```

Manages QEMU execution.

## Patch Package

### Manager

```go
type Manager struct{}

func NewManager(...) *Manager
func (pm *Manager) Apply(patchPath string) error
func (pm *Manager) List() ([]Patch, error)
```

Applies kernel patches.

## Rootfs Package

### Creator

```go
type Creator struct{}

func NewCreator(...) *Creator
func (rc *Creator) Create() error
```

Creates Debian rootfs via debootstrap.

## Toolchain Package

### Manager

```go
type Manager struct{}

func NewManager(...) *Manager
func (tm *Manager) Install() error
func (tm *Manager) Build(target string) error
func (tm *Manager) List() ([]string, error)
func (tm *Manager) Status() (map[string]bool, error)
```

Manages crosstool-ng toolchains.

All domain structs use dependency injection with infra interfaces for testability.