# Contributing to ELMOS

ELMOS is the host-side toolchain of a Master's thesis ecosystem (with
[HAVEN](https://github.com/NguyenTrongPhuc552003/haven) and
[KDAL](https://github.com/NguyenTrongPhuc552003/kdal)). Contributions are
welcome: bug reports, documentation fixes and code changes alike. Keep
changes small, reproducible and backed by a test or an observable result.

## Development setup

| Tool | Version | Check |
| ---- | ------- | ----- |
| Go | as stated in `go.mod` (currently `go 1.26.8`) | `go version` |
| Task (optional) | any v3 | `task --version` |
| QEMU (for demos and boot checks) | any recent | `qemu-system-aarch64 --version` |

The full macOS prerequisites for kernel builds are listed in the
[README Quick Start](README.md#quick-start).

```bash
git clone https://github.com/NguyenTrongPhuc552003/elmos.git
cd elmos
task build            # or: go build -o build/ ./cmd/elmos
./build/elmos doctor  # check the host environment
```

Useful [Task](https://taskfile.dev) targets (see `task --list` for all):

| Task | What it does |
| ---- | ------------ |
| `task build` | Build `build/elmos` with version information |
| `task test` | `go test ./...` |
| `task test:cover` | Tests with a coverage report (`coverage.html`) |
| `task dev:fmt` | `go fmt ./...` |
| `task dev:lint` | `golangci-lint run` if installed, else `go vet ./...` |
| `task dev:check` | fmt, lint and test in one step |
| `task dev:setup` | Download dependencies, install the complexity tool, build |

## Checks CI runs

CI ([.github/workflows/ci.yml](.github/workflows/ci.yml)) runs on every push
and pull request, on Linux, macOS and Windows. Run the same commands locally
before you open a pull request:

```bash
gofmt -l $(git ls-files '*.go')   # must print nothing (skipped on Windows in CI)
go vet ./...
go test ./...
go build -o build/ ./cmd/elmos
```

A second macOS job installs `e2fsprogs` and runs the rootfs test with
`e2fsck` available:

```bash
go test -run TestCreateDefaultsToOCIAndWritesACleanImage -v ./core/domain/rootfs/
```

ELMOS's full workflow targets macOS hosts; the Linux and Windows jobs check
that the code builds and its tests pass there.

## Commit messages

Recent history uses a bare [Conventional Commits](https://www.conventionalcommits.org/)
prefix, without a scope, and a short imperative or descriptive title:

```
feat: elmos demo linux and elmos demo haven run prebuilt demos
fix: match /init in OCI layers on Windows
docs: describe the HAVEN tasks and add haven:partb
chore: GitHub is the project home again; CI on GitHub Actions for Linux, macOS and Windows
```

Prefixes in use: `feat`, `fix`, `docs`, `chore`. Guidelines:

- One logical change per commit; keep refactors apart from behavior changes.
- The body explains why the change is needed and how it was checked (for a
  fix, name the test that fails on the old code).
- Mention user-visible changes to commands, flags or defaults.
- Use plain ASCII punctuation in messages and documentation.

## Documentation

- `docs/user/`: user guide (installation, kernel building, toolchains, QEMU).
- `docs/developer/`: developer guide (architecture, build system, testing).
- `README.md`: overview and quick start.

These folders are synced into the project documentation site
([haven-docs](https://github.com/NguyenTrongPhuc552003/haven-docs)) by its
`scripts/sync.mjs`: `README.md` becomes the ELMOS index, `docs/user` the
guide and `docs/developer` the developer section. The CLI command reference
is not written by hand: the site generates it from the Cobra definitions
with the hidden command `elmos gendocs DIR`. If you change a command, its
flags or its help text, the reference follows automatically; update the
guides when the workflow they describe changes.

## Proposing changes

1. For anything larger than a small fix, open an
   [issue](https://github.com/NguyenTrongPhuc552003/elmos/issues) first to
   discuss the problem and the approach. Bug reports and feature requests
   have templates.
2. Fork the repository and branch from `main`.
3. Make the change with tests, run the checks above, and update the docs.
4. Open a pull request against `main`. The template asks what changed, why,
   and how you tested it. CI must pass.

## Conduct and security

Please follow the [Code of Conduct](CODE_OF_CONDUCT.md). Report security
issues privately as described in [SECURITY.md](SECURITY.md), not in public
issues.

## License

By contributing, you agree that your contributions are licensed under the
[MIT License](LICENSE).
