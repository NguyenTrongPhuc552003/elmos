# Build System

ELMOS uses Task for build automation and Go's domain builders for kernel/module/app compilation.

---

## Taskfile Overview

The `Taskfile.yml` provides development tasks:

| Task               | Purpose                                    |
| ------------------ | ------------------------------------------ |
| `task build`       | Compile `build/elmos` with version ldflags |
| `task test`        | Run tests with coverage                    |
| `task dev:check`   | Format + lint                              |
| `task docs`        | Build documentation site                   |
| `task release:all` | Cross-platform Darwin builds               |

---

## Kernel Build System

### KernelBuilder

Located in `core/domain/builder/kernel.go`:

```go
type KernelBuilder struct {
    exec executor.Executor
    fs   filesystem.FileSystem
    cfg  *config.Config
    ctx  *elcontext.Context
    tm   *toolchain.Manager
}
```

### BuildOptions

```go
type BuildOptions struct {
    Jobs    int      // Parallel jobs (-j)
    Targets []string // e.g., ["Image", "dtbs", "modules"]
}
```

### Key Methods

| Method                       | Description                     |
| ---------------------------- | ------------------------------- |
| `Build(ctx, opts)`           | Execute `make` with targets     |
| `Configure(ctx, configType)` | Run menuconfig, defconfig, etc. |
| `Clean(ctx)`                 | Run `make distclean`            |
| `HasConfig()`                | Check if `.config` exists       |
| `HasKernelImage()`           | Check if kernel image built     |

### Build Flow

```
KernelBuilder.Build()
    ├── Validate targets against ValidBuildTargets
    ├── Get toolchain environment (getToolchainEnv)
    ├── Construct make arguments:
    │   - ARCH=arm64
    │   - LLVM=1
    │   - CROSS_COMPILE=<prefix>
    │   - -j<jobs>
    └── executor.RunWithEnv(make, args...)
```

---

## Module Build System

### ModuleBuilder

Located in `core/domain/builder/module.go`:

```go
type ModuleBuilder struct {
    exec executor.Executor
    fs   filesystem.FileSystem
    cfg  *config.Config
    ctx  *elcontext.Context
    tm   *toolchain.Manager
}
```

### Key Methods

| Method                   | Description                       |
| ------------------------ | --------------------------------- |
| `Build(ctx, modulePath)` | Build `.ko` from module source    |
| `Clean(ctx, modulePath)` | Clean module build artifacts      |
| `Create(name)`           | Scaffold new module from template |

---

## App Build System

### AppBuilder

Located in `core/domain/builder/app.go`:

```go
type AppBuilder struct {
    exec executor.Executor
    fs   filesystem.FileSystem
    cfg  *config.Config
    ctx  *elcontext.Context
    tm   *toolchain.Manager
}
```

Cross-compiles userspace applications for the target architecture.

---

## Environment Setup

### GetMakeEnv()

The `Context.GetMakeEnv()` method constructs the build environment:

```go
// Prepends to PATH:
// - GNU sed (libexec/gnubin)
// - GNU coreutils
// - LLVM bin
// - LLD bin
// - e2fsprogs sbin

// Sets:
// - ARCH=<target>
// - LLVM=1
// - CROSS_COMPILE=<prefix>
// - HOSTCFLAGS=<macOS compatibility flags>
```

### HOSTCFLAGS

macOS-specific flags for host tools:

```
-I<assets/libraries>     # Custom elf.h, byteswap.h
-I<libelf include>       # Homebrew libelf
-D_UUID_T
-D__GETHOSTUUID_H
-D_DARWIN_C_SOURCE
-D_FILE_OFFSET_BITS=64
```

---

## Valid Build Targets

Defined in `core/config/defaults.go`:

```go
var ValidBuildTargets = map[string]bool{
    "Image":           true,
    "zImage":          true,  // ARM32
    "dtbs":            true,
    "modules":         true,
    "modules_prepare": true,
    "all":             true,
    "vmlinux":         true,
}
```

---

## Valid Config Types

```go
var KernelConfigTypes = []string{
    "defconfig",
    "tinyconfig",
    "kvm_guest.config",
    "menuconfig",
    "xconfig",
    "nconfig",
    "oldconfig",
    "olddefconfig",
    // ...
}
```

---

## CLI Integration

### Build Command

```bash
elmos kernel build              # Default targets for arch
elmos kernel build Image        # Specific target
elmos kernel build -j 8         # Custom job count
```

### Configure Command

```bash
elmos kernel config defconfig   # Default config
elmos kernel config menuconfig  # Interactive
```

---

## Dependencies

- **Go 1.21+** - Language runtime
- **Task** - `brew install go-task`
- **LLVM** - Cross-compiler (`brew install llvm`)
- **GNU tools** - `brew install gnu-sed coreutils`