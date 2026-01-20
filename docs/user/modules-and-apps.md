# Modules and Apps

Develop kernel modules and userspace applications with ELMOS.

## Kernel Modules

### Create Module

```bash
./build/elmos module create <name>
```

Generates template in `examples/modules/<name>/` with `Makefile` and `<name>.c`.

### Build Module

```bash
cd examples/modules/<name>
make  # Uses detected toolchain
```

Outputs `.ko` file.

### Example Template

```c
#include <linux/module.h>
#include <linux/kernel.h>

static int __init hello_init(void) {
    pr_info("Hello, ELMOS!\n");
    return 0;
}

static void __exit hello_exit(void) {
    pr_info("Goodbye, ELMOS!\n");
}

module_init(hello_init);
module_exit(hello_exit);
MODULE_LICENSE("GPL");
```

### Load in QEMU

After booting kernel:

```bash
insmod /path/to/module.ko
dmesg | tail
```

## Userspace Apps

### Create App

```bash
./build/elmos app create <name>
```

Generates template in `examples/apps/<name>/` with `Makefile` and `<name>.c`.

### Build App

```bash
cd examples/apps/<name>
make  # Cross-compiles for target
```

Outputs executable.

### Example Template

```c
#include <stdio.h>

int main() {
    printf("Hello from ELMOS app!\n");
    return 0;
}
```

### Run in QEMU

Copy to rootfs and execute:

```bash
# In host
cp examples/apps/<name>/<name> build/rootfs/

# In QEMU
./<name>
```

## Cross-Compilation

- Auto-detects toolchain based on `./build/elmos arch`
- Sets `CROSS_COMPILE` and `PATH`
- Falls back to LLVM if no toolchain

## Templates

Customize templates in `assets/templates/`.

## Examples

See `examples/` for sample modules (e.g., `char-test`, `hello-world`) and apps.