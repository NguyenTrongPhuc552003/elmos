# ELMOS - Embedded Linux on MacOS

[![Build Status](https://img.shields.io/badge/build-Go%201.22+-blue)](https://github.com/NguyenTrongPhuc552003/elmos) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Native Linux kernel development tools for macOS. Build, test, and debug Linux kernels targeting RISC-V, ARM64, and more architectures—no Docker, no VMs.

## Features

- 🚀 **Native macOS builds** - Clang/LLVM toolchain with Homebrew
- 🎯 **Multi-architecture** - ARM64, RISC-V, ARM32 targets
- 📦 **Module development** - Build out-of-tree kernel modules
- 🖥️ **QEMU integration** - Run and debug kernels instantly
- 🔧 **GDB support** - Cross-architecture debugging
- 🎨 **Interactive TUI** - Menuconfig-style interface (`elmos ui`)

## Quick Start

```bash
# Install dependencies
brew install go-task

# Build from source
task build

# Check environment
./elmos doctor

# Launch interactive TUI
./elmos ui

# Or use CLI directly:
./elmos config set arch arm64
./elmos kernel config
./elmos build
./elmos qemu run
```

## Interactive TUI

Run `elmos ui` for a menuconfig-style interactive menu:

```
┌────────────────────────────────────────────────────┐
│  🔧 ELMOS - Embedded Linux on MacOS               │
├────────────────────────────────────────────────────┤
│  ▼ Setup                                           │
│      Doctor (Check Environment)              [✓]   │
│      Init Workspace                          [○]   │
│      Configure (Arch, Jobs...)                     │
│  ▼ Build                                           │
│      Build Kernel                                  │
│      Build Modules                                 │
│  ▼ Run                                             │
│      Run QEMU                                      │
├────────────────────────────────────────────────────┤
│  ↑↓: Navigate  Enter: Select  q: Quit  ?: Help    │
└────────────────────────────────────────────────────┘
```

## Build System (Task)

Uses [Task](https://taskfile.dev) instead of Make:

```bash
task --list          # Show all targets
task build           # Build elmos binary
task clean           # Clean artifacts
task deps            # Download dependencies
task fmt             # Format code
task lint            # Run linter
task test            # Run tests
task release         # Multi-platform builds
```

## Commands

| Command                          | Description                              |
| -------------------------------- | ---------------------------------------- |
| `elmos doctor`                   | Check dependencies and environment       |
| `elmos ui`                       | Launch interactive TUI menu              |
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

```bash
brew tap messense/macos-cross-toolchains
brew install llvm lld gnu-sed make libelf git qemu fakeroot e2fsprogs coreutils wget go-task
```

## Project Structure

```
.
├── apps/           # Userspace applications
├── libraries/      # macOS compatibility headers
├── modules/        # Kernel modules
├── patches/        # Kernel patches by version
├── cmd/            # CLI commands
├── internal/       # Core Go packages
│   └── tui/        # Interactive TUI (Bubbletea)
├── pkg/            # Public Go packages
├── Taskfile.yml    # Build automation (Task)
└── elmos.yaml      # Configuration (optional)
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
