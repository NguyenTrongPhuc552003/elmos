# Contributing

How to contribute to ELMOS. The repository's
[CONTRIBUTING.md](https://github.com/NguyenTrongPhuc552003/elmos/blob/main/CONTRIBUTING.md)
is the reference; this page follows it.

ELMOS is the host-side toolchain of a Master's thesis ecosystem (with
[HAVEN](https://github.com/NguyenTrongPhuc552003/haven) and
[KDAL](https://github.com/NguyenTrongPhuc552003/kdal)). Bug reports,
documentation fixes and code changes are welcome. Keep changes small,
reproducible and backed by a test or an observable result.

## Development setup

| Tool | Version | Check |
| ---- | ------- | ----- |
| Go | as stated in `go.mod` | `go version` |
| Task (optional) | any v3 | `task --version` |
| QEMU (for demos and boot checks) | any recent | `qemu-system-aarch64 --version` |

```bash
git clone https://github.com/NguyenTrongPhuc552003/elmos.git
cd elmos
task build            # or: go build -o build/ ./cmd/elmos
./build/elmos doctor  # check the host environment
```

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

`.github/workflows/ci.yml` runs on every push and pull request, on Linux,
macOS and Windows. Run the same commands before opening a pull request:

```bash
gofmt -l $(git ls-files '*.go')   # must print nothing
go vet ./...
go test ./...
go build -o build/ ./cmd/elmos
```

A second macOS job runs the rootfs test with `e2fsck` available:

```bash
go test -run TestCreateDefaultsToOCIAndWritesACleanImage -v ./core/domain/rootfs/
```

`.github/workflows/sdk.yml` builds the SDK (Linux kernel, root file system,
manifest) and publishes it as a private package for KDAL's and HAVEN's CI.
It runs on pushes to `main` that touch `sdk/**` or the workflow, or by
hand; see [SDK built by CI](sdk-ci.md). A pull request that changes those
files is checked by the run that follows its merge.

## Commit messages

A bare [Conventional Commits](https://www.conventionalcommits.org/) prefix
(`feat`, `fix`, `docs`, `chore`, `ci`), no scope required, and a short
title, for example `fix: match /init in OCI layers on Windows`.

- One logical change per commit; keep refactors apart from behavior changes.
- The body says why the change is needed and how it was checked (for a fix,
  the test that fails on the old code).
- Mention user-visible changes to commands, flags or defaults.
- Plain ASCII punctuation in messages and documentation.

## Documentation

- `docs/user/`: the user guide; `docs/developer/`: the developer guide;
  `README.md`: overview and quick start.
- The project site ([haven-docs](https://github.com/NguyenTrongPhuc552003/haven-docs))
  syncs these folders. Its CLI reference is generated from the Cobra
  definitions with the hidden command `elmos gendocs DIR`, so it follows
  the code; update the guides when a workflow they describe changes.

## Proposing changes

1. For more than a small fix, open an
   [issue](https://github.com/NguyenTrongPhuc552003/elmos/issues) first
   (bug and feature templates exist).
2. Fork, branch from `main`, make the change with tests, run the checks.
3. Open a pull request against `main`; CI must pass.

## Conduct, security, license

Follow the [Code of Conduct](https://github.com/NguyenTrongPhuc552003/elmos/blob/main/CODE_OF_CONDUCT.md).
Report security issues privately as described in
[SECURITY.md](https://github.com/NguyenTrongPhuc552003/elmos/blob/main/SECURITY.md).
Contributions are licensed under the MIT License.
