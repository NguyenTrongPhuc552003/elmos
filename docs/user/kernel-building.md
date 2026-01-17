# Kernel Building

This guide covers cloning, configuring, building, and patching Linux kernels with ELMOS.

## Supported Versions

ELMOS supports Linux v6.18+ with patches for macOS compatibility.

## Workflow

### 1. Select Architecture

```bash
./build/elmos arch <arch>  # e.g., riscv, arm64, arm
```

### 2. Clone Source

```bash
./build/elmos kernel clone
```

Clones to `build/linux/` with architecture-specific branch.

### 3. Configure

Generate default config:

```bash
./build/elmos kernel config defconfig
```

Or interactive menu:

```bash
./build/elmos kernel config menuconfig
```

### 4. Apply Patches (if needed)

For v6.18+ compatibility:

```bash
./build/elmos patch apply patches/v6.18/generic/...
```

### 5. Build

```bash
./build/elmos kernel build
```

Uses detected toolchain. Output: `build/linux/arch/<arch>/boot/Image` or `vmlinux`.

### 6. Verify

Check build artifacts:

```bash
ls build/linux/arch/*/boot/
```

## Patches

ELMOS includes patches for macOS issues:

- **v6.18**: `copy_file_range()` syscall replacement
- **ARM**: Build error fixes
- **RISC-V**: VDSO compatibility

Apply with `./build/elmos patch apply <patch>`.

## Custom Builds

- Set `HOSTCFLAGS`: ELMOS auto-sets for macOS (shims, libelf, uuid fixes)
- Environment: `CROSS_COMPILE`, `ARCH`, `PATH` set automatically
- Parallel builds: Controlled by `JOBS` config

## Troubleshooting

- "No toolchain": Install via [Toolchains](toolchains.md)
- Config errors: Regenerate defconfig
- Build hangs: Check disk space, kill and retry
- Patches fail: Ensure clean source (`./build/elmos kernel clean`)