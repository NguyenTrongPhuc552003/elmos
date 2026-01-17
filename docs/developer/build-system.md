# Build System

ELMOS uses Task for build automation and Go modules.

## Taskfile.yml

Central build configuration with tasks for:

- **build**: Compile binary with ldflags for version info
- **clean**: Remove artifacts
- **dev:fmt/lint/check**: Code quality
- **test**: Run tests with coverage
- **install**: Install binary
- **release**: Cross-platform builds

## Key Tasks

### Build

```bash
task build
```

Compiles to `build/elmos` with embedded version/commit/date.

### Development

```bash
task dev:setup  # Full setup
task dev:check  # Format and lint
task test       # Run tests
```

### Release

```bash
task release:all  # Darwin binaries + completions
task release:homebrew  # Generate Homebrew formula
```

## CI/CD

GitHub Actions (future) will run:

- `task dev:check`
- `task test`
- `task build`
- Release on tags

## Dependencies

- Go 1.21+
- Task (`brew install go-task`)
- Optional: golangci-lint, complexity tool

## Versioning

Version set via git tags, ldflags inject into `version` package.