# Installation

This guide covers installing ELMOS and its prerequisites on macOS.

## Prerequisites

ELMOS requires macOS Sequoia (15.0+) or later. Install these dependencies via Homebrew:

```bash
brew install llvm lld gnu-sed make libelf git qemu fakeroot e2fsprogs coreutils go-task wget
```

Install Xcode Command Line Tools for SDK headers:

```bash
xcode-select --install
```

!!! note
    If you encounter "gmake not found", use `brew install make` and alias `gmake` to `make`.

## Build ELMOS

Clone the repository and build the binary:

```bash
git clone https://github.com/NguyenTrongPhuc552003/elmos.git
cd elmos
task build  # Builds to build/elmos
```

## Initialize Workspace

Create a workspace with sparse image and config:

```bash
./build/elmos init
```

This generates:

- `build/elmos.sparseimage` - Workspace disk image
- `build/elmos.yaml` - Configuration file

## Verify Setup

Run the environment doctor to check dependencies:

```bash
./build/elmos doctor
```

If issues arise, see [Troubleshooting](../user/troubleshooting.md).

## Optional: Install Toolchains

For full functionality, install crosstool-ng toolchains:

```bash
./build/elmos toolchains install  # Install crosstool-ng
./build/elmos toolchains list     # List targets
./build/elmos arch riscv          # Select architecture
./build/elmos toolchains build    # Build toolchain (~30-60 min)
```

See [Toolchains](toolchains.md) for details.

## Next Steps

- [Get Started](getting-started.md) with your first kernel build
- Explore the [Interactive TUI](tui-guide.md)