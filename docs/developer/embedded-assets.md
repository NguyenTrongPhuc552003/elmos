# Embedded Assets

ELMOS embeds templates and configs using Go's `embed` package.

## Structure

```
assets/
├── embed.go          # Embed directives
├── libraries/        # Shims (elf.h, endian.h, etc.)
├── templates/
│   ├── app/          # App templates (main.c, Makefile)
│   ├── configs/      # YAML configs (elmos.yaml)
│   ├── init/         # Init scripts (guesync.sh, init.sh)
│   └── module/       # Module templates (module.c, Makefile)
└── toolchains/
    └── configs/      # Crosstool-ng configs
```

## Embedding

In `assets/embed.go`:

```go
//go:embed libraries/*
//go:embed templates/*
//go:embed toolchains/configs/*
var FS embed.FS
```

Accessed via `assets.FS`.

## Usage

Templates used for code generation:

- **Modules**: `templates/module/` for `elmos module create`
- **Apps**: `templates/app/` for `elmos app create`
- **Configs**: `templates/configs/` for workspace init

## Libraries

Shims for macOS compatibility:

- `elf.h`: libelf headers
- `endian.h`: Byte order functions
- `asm/bitsperlong.h`: Architecture defines

Used in kernel builds to replace missing macOS headers.

## Toolchains

Pre-configured crosstool-ng `.config` files for targets.

## Maintenance

Update templates in `assets/`, rebuild to embed changes.