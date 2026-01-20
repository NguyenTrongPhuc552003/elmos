# Embedded Assets

ELMOS embeds templates and headers at compile time using Go's `embed` package.

---

## Directory Structure

```
assets/
├── embed.go              # Embed directives and accessor functions
├── libraries/            # macOS compatibility headers
│   ├── elf.h             # ELF definitions
│   ├── byteswap.h        # Byte swapping macros
│   └── asm/
│       └── bitsperlong.h # Architecture bit width
├── templates/
│   ├── app/              # Userspace app templates
│   │   ├── main.c.tmpl
│   │   └── Makefile.tmpl
│   ├── configs/          # Configuration templates
│   │   └── elmos.yaml.tmpl
│   ├── init/             # Guest init scripts
│   │   ├── init.sh.tmpl
│   │   └── guesync.sh.tmpl
│   └── module/           # Kernel module templates
│       ├── module.c.tmpl
│       └── Makefile.tmpl
└── toolchains/
    └── configs/          # Crosstool-ng configurations
```

---

## Embed Directives

```go
// assets/embed.go
package assets

import "embed"

//go:embed templates/*
var Templates embed.FS
```

---

## Accessor Functions

| Function              | Returns                             |
| --------------------- | ----------------------------------- |
| `GetModuleTemplate()` | `templates/module/module.c.tmpl`    |
| `GetModuleMakefile()` | `templates/module/Makefile.tmpl`    |
| `GetAppTemplate()`    | `templates/app/main.c.tmpl`         |
| `GetAppMakefile()`    | `templates/app/Makefile.tmpl`       |
| `GetInitScript()`     | `templates/init/init.sh.tmpl`       |
| `GetGuestSync()`      | `templates/init/guesync.sh.tmpl`    |
| `GetConfigTemplate()` | `templates/configs/elmos.yaml.tmpl` |

**Usage:**

```go
tmpl, err := assets.GetModuleTemplate()
if err != nil {
    return err
}
// Use tmpl bytes...
```

---

## Template Variables

### Module Template

```c
// templates/module/module.c.tmpl
#include <linux/module.h>
#include <linux/kernel.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("{{.Author}}");
MODULE_DESCRIPTION("{{.Description}}");

static int __init {{.Name}}_init(void) {
    pr_info("{{.Name}}: loaded\n");
    return 0;
}
module_init({{.Name}}_init);
```

### App Template

```c
// templates/app/main.c.tmpl
#include <stdio.h>

int main(void) {
    printf("Hello from {{.Name}}!\n");
    return 0;
}
```

---

## macOS Compatibility Headers

The `assets/libraries/` directory contains headers missing on macOS:

| Header              | Purpose                |
| ------------------- | ---------------------- |
| `elf.h`             | ELF format definitions |
| `byteswap.h`        | Byte swapping macros   |
| `asm/bitsperlong.h` | Architecture bit width |

These are included via `HOSTCFLAGS=-I<assets/libraries>`.

---

## Adding New Templates

1. Create template in `assets/templates/<category>/`
2. Add accessor function in `assets/embed.go`:
   ```go
   func GetNewTemplate() ([]byte, error) {
       return Templates.ReadFile("templates/category/new.tmpl")
   }
   ```
3. Rebuild: `task build`