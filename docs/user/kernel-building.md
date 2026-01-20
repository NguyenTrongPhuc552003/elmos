# Kernel Building

Build Linux kernels for ARM64, ARM, and RISC-V on macOS.

---

## Prerequisites

- Workspace initialized: `elmos init`
- Dependencies checked: `elmos doctor`

---

## Build Workflow

### 1. Set Architecture

```bash
elmos arch arm64    # or: arm, riscv
elmos arch          # Show current
```

### 2. Configure Kernel

```bash
# Default config for architecture
elmos kernel config defconfig

# Interactive menu
elmos kernel config menuconfig

# Minimal config
elmos kernel config tinyconfig
```

**Valid config types:**

| Type               | Description          |
| ------------------ | -------------------- |
| `defconfig`        | Architecture default |
| `tinyconfig`       | Minimal kernel       |
| `menuconfig`       | Interactive menu     |
| `kvm_guest.config` | KVM optimized        |
| `oldconfig`        | Update existing      |
| `olddefconfig`     | Update with defaults |

### 3. Build

```bash
# Default targets (Image, dtbs, modules)
elmos kernel build

# Specific targets
elmos kernel build Image
elmos kernel build vmlinux

# Custom parallelism
elmos kernel build -j 8
```

**Valid build targets:**

| Target    | Description                   |
| --------- | ----------------------------- |
| `Image`   | Kernel image (arm64, riscv)   |
| `zImage`  | Compressed image (arm)        |
| `dtbs`    | Device tree blobs             |
| `modules` | Kernel modules                |
| `vmlinux` | Uncompressed kernel (for GDB) |

### 4. Verify Build

```bash
elmos status
```

Output:

```
Workspace Status:
  Volume: /Volumes/elmos (mounted)
  Kernel: ✓ Configured, ✓ Built
  Architecture: arm64
  Image: /Volumes/elmos/linux/arch/arm64/boot/Image
```

---

## BuildOptions Reference

```go
// core/domain/builder/kernel.go
type BuildOptions struct {
    Jobs    int      // Parallel jobs (-j)
    Targets []string // Build targets
}
```

---

## Environment Variables

ELMOS automatically sets:

| Variable        | Value                     |
| --------------- | ------------------------- |
| `ARCH`          | Target architecture       |
| `LLVM`          | `1` (use LLVM toolchain)  |
| `CROSS_COMPILE` | Toolchain prefix          |
| `HOSTCFLAGS`    | macOS compatibility flags |
| `PATH`          | Prepends LLVM, GNU tools  |

---

## Clean Build

```bash
elmos kernel clean    # make distclean
```

---

## Patches

Apply macOS compatibility patches:

```bash
# List available
elmos patch list

# Apply
elmos patch apply v6.18/generic/fix-copy-range
```

---

## Troubleshooting

| Issue           | Solution                                  |
| --------------- | ----------------------------------------- |
| "No toolchain"  | Run `elmos doctor`                        |
| Config errors   | Run `elmos kernel clean` then reconfigure |
| Build hangs     | Check disk space, reduce `-j`             |
| Missing headers | Run `elmos doctor --fix`                  |