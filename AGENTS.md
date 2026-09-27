# ELMOS agent guidance

ELMOS is a Go CLI for building and testing Embedded Linux projects on macOS. Keep investigations scoped to the user's task. This file is a routing guide, not a requirement to read the whole repository before every change.

## Where to look

- `cmd/elmos/main.go` starts the CLI; `core/app/app.go` wires dependencies and the Cobra root command; `core/app/commands/` registers user commands.
- `core/config/` loads and saves configuration; `core/context/` carries runtime context.
- `core/domain/` implements builders, toolchain management, doctor checks, patches, rootfs creation, and QEMU behavior.
- `core/infra/` provides shell execution, filesystem access, and Homebrew resolution; `core/ui/` handles terminal output and TUI.
- `assets/embed.go`, `assets/templates/`, and `assets/toolchains/configs/` supply embedded resources. `patches/` contains versioned Linux patches. `examples/` contains sample apps and modules.
- `Taskfile.yml` defines development commands. Read `docs/developer/architecture.md` for architecture questions, `docs/developer/build-system.md` for task/build questions, and the relevant `docs/user/` page for user-visible behavior. Verify claims against code when docs disagree.

## Working method

1. Check `git status --short` before edits. Use `git ls-files` to locate project files, then `rg -n` for symbols or behavior. Open only the matching files and their direct callers or dependencies. Expand when a concrete question remains.
2. Trace a CLI behavior from `cmd/elmos/main.go` through `core/app/commands/` to the domain and infrastructure functions actually invoked. Distinguish observed code behavior from documentation or assumptions.
3. Keep shell output bounded. Prefer targeted paths and matches over recursive dumps. Do not read an entire Linux source tree, toolchain tree, generated image, or build log unless the task specifically requires it.
4. Treat `build/`, `data/`, sparse images, generated binaries, cloned kernel sources, and installed toolchains as workspace outputs or external dependencies. Inspect them only when relevant. The symlinks in `assets/libraries/asm/` may point into a local Linux checkout; do not follow them during ordinary repository exploration.
5. Make the smallest relevant change, follow existing package boundaries, and preserve unrelated work. For code changes, run focused Go tests or checks first; broaden verification only when the change warrants it. Report commands run and remaining limitations.

## Side effects

Exploration is read only. Do not run setup, installation, toolchain builds, kernel cloning/building, rootfs creation, QEMU, `task clean`, or `task dev:check` merely to understand the project. Some tasks modify files, download dependencies, or consume substantial disk and time. Use them when the user requests the related work and the effect is understood. `task dev:fmt` and `task dev:deps` can modify source or module files; inspect their definitions before using them.
