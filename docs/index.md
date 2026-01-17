# ELMOS Documentation

Welcome to the official documentation for **ELMOS** (Embedded Linux on MacOS), a complete embedded Linux SDK for macOS. ELMOS enables native kernel development, cross-compilation, and emulation without Docker or VMs.

## Quick Links

- [Get Started](user/getting-started.md) - Your first kernel build
- [Installation](user/installation.md) - Setup and prerequisites
- [Contributing](developer/contributing.md) - Help improve ELMOS

## Overview

ELMOS provides:

- **Native Toolchain Management**: Build and manage cross-compilers for ARM64, ARM, and RISC-V using crosstool-ng.
- **Kernel Automation**: Clone, configure, build, and test Linux kernels (v6.18+).
- **Interactive TUI**: Rich terminal interface for streamlined workflows.
- **QEMU Integration**: Boot and debug kernels with GDB support.
- **Module & App Development**: Cross-compile kernel modules and userspace applications.

Built with Go, ELMOS leverages macOS tools like LLVM and Homebrew for seamless development.

## Documentation Sections

### User Guide
For users installing and using ELMOS:

- Installation and setup
- Tutorials and examples
- Command references
- Troubleshooting

### Developer Guide
For contributors:

- Architecture and design
- API documentation
- Contributing guidelines
- Testing and build processes

## Support

- [GitHub Issues](https://github.com/NguyenTrongPhuc552003/elmos/issues) for bugs and features
- [Discussions](https://github.com/NguyenTrongPhuc552003/elmos/discussions) for questions

---

*ELMOS is MIT-licensed. Inspired by [Seiya's tutorial](https://seiya.me/blog/building-linux-on-macos-natively).*"