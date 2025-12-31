# ELMOS - Embedded Linux on MacOS

[![Build Status](https://img.shields.io/badge/build-Go%201.22+-blue)](https://github.com/NguyenTrongPhuc552003/elmos) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Native Linux kernel development tools for macOS. Build, test, and debug Linux kernels targeting RISC-V, ARM64, and more architectures—no Docker, no VMs.

## Features

- 🚀 **Native macOS builds** - Clang/LLVM toolchain with Homebrew
- 🎯 **Multi-architecture** - ARM64, RISC-V, ARM32 targets
- 📦 **Module development** - Build out-of-tree kernel modules
- 🖥️ **QEMU integration** - Run and debug kernels instantly
- 🔧 **GDB support** - Cross-architecture debugging
- 🎨 **Modern CLI** - Cobra-based with styled output

## Quick Start

```bash
# Install (from source)
make build
make install

# Check environment
elmos doctor

# Initialize workspace
elmos init

# Configure and build for ARM64
elmos config set arch arm64
elmos kernel config
elmos build

# Create rootfs and run in QEMU
elmos rootfs create
elmos qemu run
```

## Commands

| Command                          | Description                              |
| -------------------------------- | ---------------------------------------- |
| `elmos doctor`                   | Check dependencies and environment       |
| `elmos init`                     | Mount workspace and clone kernel         |
| `elmos image mount/unmount`      | Manage sparse disk image                 |
| `elmos repo checkout <tag>`      | Checkout kernel version                  |
| `elmos config set <key> <value>` | Configure settings                       |
| `elmos kernel config`            | Configure kernel (defconfig, menuconfig) |
| `elmos build`                    | Build kernel, dtbs, modules              |
| `elmos module build [name]`      | Build kernel modules                     |
| `elmos app build [name]`         | Build userspace apps                     |
| `elmos rootfs create`            | Create Debian rootfs                     |
| `elmos qemu run`                 | Run kernel in QEMU                       |
| `elmos qemu debug`               | Run with GDB stub                        |
| `elmos patch apply <file>`       | Apply kernel patches                     |

## Configuration

Create `elmos.yaml` in your project root:

```yaml
image:
  volume_name: kernel-dev
  size: 20G

build:
  arch: arm64
  jobs: 8
  llvm: true

qemu:
  memory: 2G
  gdb_port: 1234

profiles:
  riscv-dev:
    arch: riscv
    memory: 2G
  arm64-dev:
    arch: arm64
    memory: 4G
```

Apply profiles: `elmos config profile riscv-dev`

## Dependencies

Install via Homebrew:

```bash
brew tap messense/macos-cross-toolchains
brew install llvm lld gnu-sed make libelf git qemu fakeroot e2fsprogs coreutils wget
```

## Project Structure

```
.
├── apps/           # Userspace applications
├── libraries/      # macOS compatibility headers (elf.h, byteswap.h)
├── modules/        # Kernel modules
│   ├── hello-world/
│   └── ...
├── patches/        # Kernel patches by version
│   └── v6.18/
├── cmd/            # CLI commands
├── internal/       # Core Go packages
├── pkg/            # Public Go packages
├── go.mod          # Go module
├── Makefile        # Build automation
└── elmos.yaml      # Configuration (optional)
```

## Kernel Modules

```bash
# Create new module
elmos module new my-driver

# Build modules
elmos module build my-driver

# Check status
elmos module status

# Prepare headers for building
elmos module headers
```

## Userspace Apps

```bash
# Create new app
elmos app new my-app

# Build for target architecture
elmos app build my-app

# List available apps
elmos app list
```

## Debugging

```bash
# Terminal 1: Start QEMU with GDB stub
elmos qemu debug

# Terminal 2: Connect GDB
elmos qemu gdb
```

## License

MIT - See [LICENSE](LICENSE)

## Author

Phuc Nguyen ([@NguyenTrongPhuc552003](https://github.com/NguyenTrongPhuc552003))
