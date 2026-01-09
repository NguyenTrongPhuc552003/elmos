# ELMOS – Embedded Linux on MacOS

[![Build Status](https://img.shields.io/badge/build-v6.18%20ARM64-green)](https://github.com/NguyenTrongPhuc552003/elmos) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A complete embedded Linux SDK for macOS. Build kernels, cross-compile with native toolchains, develop kernel modules and userspace apps—all without Docker or VMs. ELMOS provides an integrated development environment with interactive TUI, automatic toolchain management (crosstool-ng), and seamless QEMU integration. Targeting RISC-V, ARM64, ARM, and more, with full support for Linux v6.18+.

Inspired by [Seiya's tutorial](https://seiya.me/blog/building-linux-on-macos-natively) (which fixed older kernels like v6.17), we extended it for v6.18's new challenges, like the `copy_file_range()` incompatibility.

## Features

- **Cross-Compiler Toolchain Management**: Build and manage crosstool-ng toolchains for ARM64, ARM, and RISC-V
- **Interactive TUI**: Rich terminal interface with real-time command output
- **Environment Doctor**: Comprehensive dependency and toolchain health checks
- **Kernel Build Automation**: Configure, build, and test Linux kernels with integrated toolchain support
- **Module & App Development**: Build kernel modules and userspace apps with automatic cross-compilation
- **QEMU Integration**: Boot and debug kernels with built-in GDB support

## Quick Start

### 1. Prerequisites (macOS Sequoia/Tahoe+ recommended)

```bash
brew install llvm lld gnu-sed make libelf git qemu fakeroot e2fsprogs coreutils go-task wget
xcode-select --install  # For SDK headers
```

### 2. Build & Initialize

```bash
git clone https://github.com/NguyenTrongPhuc552003/elmos.git
cd elmos
task build                # Build to build/elmos
./build/elmos init        # Create workspace (sparseimage + config)
./build/elmos doctor      # Verify environment
```

### 3. Install Toolchain (Optional but Recommended)

```bash
./build/elmos toolchains install       # Install crosstool-ng
./build/elmos toolchains list          # List available targets
./build/elmos arch riscv               # Select arch (auto-selects toolchain)
./build/elmos toolchains build         # Build the toolchain (~30-60 min)
./build/elmos toolchains status        # Verify installation
```

### 4. Configure & Build Kernel

```bash
./build/elmos kernel clone             # Clone kernel source
./build/elmos kernel config defconfig  # Or: menuconfig
./build/elmos kernel build             # Build with detected toolchain
```

### 5. Create RootFS & Run

```bash
./build/elmos rootfs create            # Debian rootfs (debootstrap)
./build/elmos qemu run                 # Boot in QEMU
./build/elmos qemu debug               # With GDB stub (port 1234)
```

## Interactive TUI

Launch with `./build/elmos tui` for a rich interactive interface:

![elmos TUI](docs/images/elmos_tui.png)

## Toolchain Management

ELMOS integrates [crosstool-ng](https://crosstool-ng.github.io/) for building native cross-compilers:

| Command                       | Description                                         |
| ----------------------------- | --------------------------------------------------- |
| `elmos toolchains install`    | Clone & build crosstool-ng                          |
| `elmos toolchains list`       | List available target configurations                |
| `elmos toolchains <target>`   | Select a target (e.g., `riscv64-unknown-linux-gnu`) |
| `elmos toolchains build`      | Build selected toolchain                            |
| `elmos toolchains status`     | Show installed toolchains                           |
| `elmos toolchains env`        | Display environment variables                       |
| `elmos toolchains menuconfig` | Interactive toolchain configuration                 |
| `elmos toolchains clean`      | Clean toolchain build artifacts                     |

**Pre-configured targets** with optimized settings:
- `aarch64-unknown-linux-gnu` (ARM64)
- `arm-cortex_a15-linux-gnueabihf` (ARM 32-bit)
- `riscv64-unknown-linux-gnu` (RISC-V 64-bit)

## Build System (Task)

Uses [Task](https://taskfile.dev) with namespaced commands:

```bash
task --list              # Show all targets

# Core
task build               # Build elmos binary → build/elmos
task clean               # Clean all artifacts

# Development
task dev:check           # Run fmt, lint, test (pre-commit style)
task dev:setup           # Full setup (deps + build + init)
task test                # Run tests
task test:cover          # Tests with coverage report

# elmos CLI Wrappers
task elmos:init          # Initialize workspace
task elmos:doctor        # Run environment check
task elmos:status        # Show workspace status
task elmos:tui           # Launch interactive TUI

# Release
task release:darwin      # Build for macOS (arm64 + amd64)
task release:all         # Full release with completions
```

## Repo Structure

```bash
.
├── assets/               # Embedded templates
│   └── templates/        # Config, module, app templates
├── build/                # Build output (elmos, elmos.yaml, img.sparseimage)
├── core/                 # Core domain logic
│   ├── app/              # CLI application & command wiring
│   ├── config/           # Configuration management
│   ├── domain/           # Business logic (builder, toolchain, doctor)
│   └── ui/               # User Interface (TUI)
├── libraries/            # Shims: byteswap.h, elf.h, asm/
├── modules/              # Kernel modules
├── patches/              # Versioned kernel patches
├── tools/
│   └── toolchains/       # crosstool-ng & custom configs
├── Taskfile.yml          # Build automation
└── main.go               # Entry point
```

## Key Workarounds

### 1. v6.18 `copy_file_range()` Incompatibility
- **Issue**: v6.18 uses Linux-only syscall in `gen_init_cpio`
- **Fix**: Patch replaces with `copyfile(COPYFILE_DATA)` on macOS
- **Apply**: `./build/elmos patch apply patches/v6.18/0001-usr-gen_init_cpio-Replace-linux-kernel-syscall-with-.patch`

### 2. Automatic Toolchain Detection
- Kernel, module, and app builds auto-detect installed toolchains
- `CROSS_COMPILE` and `PATH` set automatically for the target architecture
- Falls back to Homebrew LLVM if no toolchain installed

### 3. HOSTCFLAGS for macOS
The CLI sets these automatically:
- `-I${MACOS_HEADERS}`: Custom shims (`elf.h`, `byteswap.h`)
- `-I${LIBELF_INCLUDE}`: libelf for ELF parsing
- `-D_UUID_T -D__GETHOSTUUID_H`: Suppress uuid_t conflicts
- `-D_DARWIN_C_SOURCE`: macOS 10.15+ APIs

## Troubleshooting

| Issue                 | Solution                              |
| --------------------- | ------------------------------------- |
| "gmake not found"     | `brew install make` → use `gmake`     |
| UUID conflicts        | Ensure patch applied                  |
| Toolchain build fails | Check `elmos doctor` for missing deps |
| TUI shows help text   | Rebuild with `task build`             |

## Credits

- **Original Tutorial**: [Building Linux on macOS Natively](https://seiya.me/blog/building-linux-on-macos-natively) by Seiya Suzuki
- **Upstream**: [Clang Built Linux](https://clangbuiltlinux.github.io/) for LLVM guidance
- **Author**: Phuc Nguyen ([@NguyenTrongPhuc552003](https://github.com/NguyenTrongPhuc552003))

## License

MIT — fork, extend, build freely.
