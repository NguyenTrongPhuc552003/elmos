# Getting Started

This tutorial walks through building your first Linux kernel with ELMOS on macOS.

## Prerequisites

Ensure ELMOS is [installed](installation.md) and the workspace initialized.

## Step 1: Select Architecture

Choose a target architecture (e.g., RISC-V):

```bash
./build/elmos arch riscv
```

This auto-selects the `riscv64-unknown-linux-gnu` toolchain if installed.

## Step 2: Clone Kernel Source

Clone the Linux kernel repository:

```bash
./build/elmos kernel clone
```

This makes a shallow clone of the pinned tag (v6.18) into the workspace
volume and applies the macOS host patches that elmos ships for that tag.
Use `--ref <tag>` for another tag and `--full` for the whole history.

## Step 3: Configure Kernel

Generate a default config:

```bash
./build/elmos kernel config defconfig
```

For custom config, use menuconfig:

```bash
./build/elmos kernel config menuconfig
```

## Step 4: Build Kernel

Build the kernel with the detected toolchain:

```bash
./build/elmos kernel build
```

This may take 10-30 minutes depending on hardware.

## Step 5: Create RootFS

Create the root filesystem:

```bash
./build/elmos rootfs create
```

By default this pulls `alpine:3` for the target architecture and writes it
straight into an ext4 image (a few seconds, no sudo). Alpine's busybox
provides `insmod` and `ip`, so modules load and the network comes up.
Choose another image with `--image debian:stable-slim` (no `insmod`, no
`ip`: add them yourself), or the debootstrap path with
`--engine debootstrap` (needs sudo).

## Step 6: Run in QEMU

Boot the kernel in QEMU:

```bash
./build/elmos qemu run
```

You should see the Linux boot process. Login with `root` (no password).

## Step 7: Debug (Optional)

For debugging, run with GDB stub:

```bash
./build/elmos qemu debug
```

Connect GDB in another terminal:

```bash
gdb-multiarch build/linux/vmlinux
(gdb) target remote :1234
```

## Next Steps

- [Develop modules](modules-and-apps.md)
- [Customize toolchains](toolchains.md)
- [Use the TUI](tui-guide.md) for interactive workflows

## Troubleshooting

If builds fail, check [Troubleshooting](troubleshooting.md) or run `./build/elmos doctor`.