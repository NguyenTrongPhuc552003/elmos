# QEMU Integration

Run and debug kernels using ELMOS's QEMU integration.

## Running Kernels

### Basic Run

```bash
./build/elmos qemu run
```

Boots kernel with generated rootfs. Requires kernel and rootfs built.

### Debug Mode

```bash
./build/elmos qemu debug
```

Starts QEMU with GDB stub on port 1234.

Connect GDB:

```bash
gdb-multiarch build/linux/vmlinux
(gdb) target remote :1234
(gdb) continue
```

### Custom Options

Modify `elmos.yaml` or use `./build/elmos qemu options` to set QEMU args.

## Supported Architectures

- **RISC-V**: `qemu-system-riscv64`
- **ARM64**: `qemu-system-aarch64`
- **ARM**: `qemu-system-arm`

## Networking

QEMU runs with user-mode networking. Access host via `10.0.2.2`.

## Disk Images

Uses `build/elmos.sparseimage` for workspace. Rootfs mounted as `/dev/vda`.

## Troubleshooting

- "Kernel not found": Build kernel first
- "No rootfs": Create with `./build/elmos rootfs create`
- Boot hangs: Check kernel config for console/serial
- GDB fails: Ensure `gdb-multiarch` installed (`brew install gdb`)