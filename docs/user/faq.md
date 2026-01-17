# FAQ

Frequently asked questions about ELMOS.

## General

**What is ELMOS?**  
A native SDK for building Linux kernels on macOS without VMs.

**Why macOS?**  
Leverages LLVM, Homebrew, and native tools for seamless development.

**Supported architectures?**  
ARM64, ARM, RISC-V (Linux v6.18+).

## Installation

**Do I need Docker?**  
No, ELMOS is fully native.

**Can I use older macOS?**  
Requires Sequoia+ for compatibility.

## Toolchains

**Are toolchains required?**  
Optional; falls back to LLVM, but recommended for full features.

**How long to build toolchain?**  
30-60 min, depending on hardware.

## Kernel Building

**Which kernel versions?**  
v6.18+ with patches for macOS.

**Can I use custom configs?**  
Yes, via menuconfig.

## Development

**How to develop modules/apps?**  
Use `./build/elmos module/app create`, then `make`.

**Cross-compilation?**  
Automatic with detected toolchains.

## QEMU

**Networking in QEMU?**  
User-mode; host at 10.0.2.2.

**Debugging?**  
Use `./build/elmos qemu debug` with GDB.

## Contributing

**How to contribute?**  
See [Developer Guide](../developer/contributing.md).

**License?**  
MIT.