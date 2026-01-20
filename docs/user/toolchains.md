# Toolchain Management

ELMOS integrates [crosstool-ng](https://crosstool-ng.github.io/) for building native cross-compilers on macOS.

## Overview

Toolchains enable cross-compilation for target architectures without VMs. ELMOS supports pre-configured targets with optimized settings.

## Supported Targets

| Target                           | Architecture | Description             |
| -------------------------------- | ------------ | ----------------------- |
| `aarch64-unknown-linux-gnu`      | ARM64        | 64-bit ARM              |
| `arm-cortex_a15-linux-gnueabihf` | ARM          | 32-bit ARM (Cortex-A15) |
| `riscv64-unknown-linux-gnu`      | RISC-V       | 64-bit RISC-V           |

## Commands

### Install crosstool-ng

```bash
./build/elmos toolchains install
```

Clones and builds crosstool-ng to `build/toolchains/crosstool-ng/`.

### List Targets

```bash
./build/elmos toolchains list
```

Shows available configurations.

### Select Target

```bash
./build/elmos toolchains <target>
```

Example: `./build/elmos toolchains riscv64-unknown-linux-gnu`

### Build Toolchain

```bash
./build/elmos toolchains build
```

Builds the selected toolchain (~30-60 min). Installs to `build/toolchains/install/`.

### Check Status

```bash
./build/elmos toolchains status
```

Verifies installation and environment variables.

### Show Environment

```bash
./build/elmos toolchains env
```

Displays `CROSS_COMPILE`, `PATH`, etc.

### Customize Config

```bash
./build/elmos toolchains menuconfig
```

Interactive configuration for advanced users.

### Clean Artifacts

```bash
./build/elmos toolchains clean
```

Removes build artifacts.

## Automatic Detection

Kernel, module, and app builds auto-detect installed toolchains based on selected architecture (`./build/elmos arch <arch>`).

Falls back to Homebrew LLVM if no toolchain installed.

## Custom Toolchains

For custom targets, modify configs in `assets/toolchains/configs/` and rebuild.

## Troubleshooting

- Build fails: Ensure all deps installed (`./build/elmos doctor`)
- Slow builds: Use more cores with `CT_PARALLEL_JOBS` env var
- Conflicts: Clean and rebuild