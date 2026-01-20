# User Guide Overview

This section provides comprehensive guides for using ELMOS to build and develop Linux kernels on macOS.

## Getting Started

If you're new to ELMOS:

1. [Install ELMOS](installation.md) - Prerequisites and setup
2. [First Kernel Build](getting-started.md) - Step-by-step tutorial
3. [Toolchain Management](toolchains.md) - Install and configure cross-compilers

## Core Workflows

- **Kernel Development**: [Clone, configure, and build kernels](kernel-building.md)
- **Module & App Creation**: [Develop kernel modules and userspace apps](modules-and-apps.md)
- **Emulation**: [Run and debug with QEMU](qemu-integration.md)
- **Interactive Mode**: [Use the TUI](tui-guide.md)

## Reference

- [Troubleshooting](troubleshooting.md) - Common issues and solutions
- [FAQ](faq.md) - Frequently asked questions
- [Changelog](changelog.md) - What's new

## Prerequisites

ELMOS requires macOS Sequoia or later with:

- Homebrew for dependencies
- Xcode Command Line Tools
- Basic familiarity with terminal commands

For advanced features, install crosstool-ng toolchains.

## Support

Encounter an issue? Check [Troubleshooting](troubleshooting.md) or open a [GitHub Issue](https://github.com/NguyenTrongPhuc552003/elmos/issues).