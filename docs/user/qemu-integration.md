# QEMU Integration

Run and debug kernels using ELMOS's QEMU integration.

---

## CLI Reference

```bash
elmos qemu [flags]

Flags:
  -r, --run         Run mode (boot kernel)
  -d, --debug       Debug mode (tmux + GDB)
  -t, --target      Path to app or module (repeatable)
  -l, --list        List available machines
  -p, --pick        Select specific machine
      --graphical   Use graphical display
```

---

## Run Modes

### Basic Run

```bash
elmos qemu -r
```

Boots kernel with generated rootfs. Output:

```
→ Starting QEMU...
[    0.000000] Booting Linux on physical CPU 0x0000000000
[    0.000000] Linux version 6.0.0-dirty ...
...
System ready.
#
```

### Debug Mode

```bash
elmos qemu -d
```

Opens tmux session with:

- **Left pane**: QEMU (paused at boot)
- **Right pane**: GDB connected to kernel

```bash
# In GDB pane
(gdb) break start_kernel
(gdb) continue
```

### With Targets

Load userspace apps or kernel modules:

```bash
# Run with user application
elmos qemu -r -t ./examples/apps/hello/hello

# Debug with kernel module
elmos qemu -d -t ./examples/modules/hello/hello.ko
```

Targets are synced to `/mnt/share` inside the guest.

---

## Machine Selection

### List Machines

```bash
elmos qemu -l
```

Shows available QEMU machines for current architecture:

```
ℹ Available QEMU Machines for arm64:
  * virt - QEMU 10.2 ARM Virtual Machine (default)
    raspi3b - Raspberry Pi 3B
    raspi4b - Raspberry Pi 4B
    sbsa-ref - QEMU SBSA Reference
```

### Pick Machine

```bash
elmos qemu -p raspi4b -r
```

Uses Raspberry Pi 4B machine instead of default `virt`.

---

## RunOptions (Developer Reference)

```go
// core/domain/emulator/options.go
type RunOptions struct {
    Debug     bool     // Enable GDB stub
    Run       bool     // Run mode
    Graphical bool     // GUI display
    Targets   []Target // Apps/modules to load
    Machine   string   // Override machine
}
```

---

## Architecture Defaults

| Arch  | QEMU Binary           | Default Machine    | Console   |
| ----- | --------------------- | ------------------ | --------- |
| arm64 | `qemu-system-aarch64` | `virt`             | `ttyAMA0` |
| arm   | `qemu-system-arm`     | `virt,highmem=off` | `ttyAMA0` |
| riscv | `qemu-system-riscv64` | `virt`             | `ttyS0`   |

---

## Graphical Mode

```bash
elmos qemu -r --graphical
```

Opens QEMU with GUI window (requires virtio-gpu kernel config).

---

## Networking

QEMU runs with user-mode networking:

- Guest can access internet
- Host accessible at `10.0.2.2`
- SSH forwarded: host `:2222` → guest `:22`

```bash
# From host
ssh -p 2222 root@localhost
```

---

## Troubleshooting

| Issue              | Solution                                  |
| ------------------ | ----------------------------------------- |
| "Kernel not found" | Run `elmos kernel build` first            |
| "No rootfs"        | Run `elmos rootfs create`                 |
| Boot hangs         | Check kernel config for `CONFIG_SERIAL_*` |
| Invalid machine    | Run `elmos qemu -l` to see valid options  |
| GDB fails          | Install `gdb` via Homebrew                |