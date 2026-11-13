# Changelog

All notable changes to this project are documented in this file.
The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Release notes for each version are on its GitHub release page.

## [Unreleased]

Changes since v3.2.0, summarised from the commit history.

### Added

- `elmos demo linux` and `elmos demo haven` run prebuilt demos with only
  QEMU: no kernel build, no workspace volume, no sudo. `demo linux` boots
  Linux v6.18 arm64 with an alpine:3 root file system and loads KDAL;
  `demo haven` runs the HAVEN try bundle. Assets are checked against the
  release's SHA256SUMS and cached in `~/.elmos/demo`, so later runs work
  offline. A release not yet on GitHub is fetched from its Codeberg copy.
- OCI rootfs engine, now the default for `rootfs create`: it pulls a
  container image for the target architecture and writes ext4 without
  extracting anything onto the Mac, so ownership, setuid bits and
  case-sensitive names survive without sudo or fakeroot. debootstrap
  remains available as `--engine debootstrap`.
- Hidden `elmos gendocs DIR` command: writes one Markdown page per command
  (usage, flags, subcommand links) from the Cobra definitions, used by the
  documentation site so the CLI reference cannot drift from the code.
- Kernel patches and the Linux headers macOS lacks are embedded in the
  binary; `kernel clone` applies the patches for its tag and `patch list`
  shows them. `init`, `kernel build` and `doctor --fix` write the headers
  into the project without network access.
- Taskfile tasks to build, test and boot HAVEN from ELMOS: `haven:test`,
  `haven:build`, `haven:smoke`, `haven:partb` and `haven:all`; README
  section on the thesis toolchain.
- CI on GitHub Actions: gofmt check, `go vet`, `go test` and `go build` on
  Linux, macOS and Windows, plus an `e2fsck` check of the rootfs image on
  macOS.
- SDK built by CI (`.github/workflows/sdk.yml`, `docs/developer/sdk-ci.md`):
  Linux v6.18 arm64 (defconfig plus KUnit, LLVM, neutral build names), an
  alpine:3 root file system and a manifest, smoke-booted in QEMU and
  published as the private package `ghcr.io/nguyentrongphuc552003/elmos-sdk`
  for KDAL's and HAVEN's CI.
- `executor.Recorder`, a host-free executor for domain tests.
- API reference overview page in the developer documentation.

### Changed

- The OCI engine writes ext4 in pure Go (an adapted copy of hcsshim's ext4
  writer in `third_party/ext4`). e2fsprogs is now optional (debootstrap
  engine only) and dropping the hcsshim module reduces the binary from
  22.0 MB to 19.3 MB.
- `kernel clone` and `kernel reset` default to a shallow clone of v6.18
  (`--ref`, `--full`, `--no-patches` to change this).
- The guest `/init` is embedded and shared by both rootfs engines instead
  of being read from a repository checkout.
- Documentation marks the native macOS kernel build as experimental and
  describes cloning, rootfs engines, embedded patches and headers.
- Project home moved to Codeberg and then back to GitHub; the Go module
  path is `github.com/NguyenTrongPhuc552003/elmos` again, and links, the
  demo release base and the Homebrew formula point at GitHub.
- Removed the `docs.yml` GitHub Pages workflow and `AGENTS.md`.
- Em-dashes in the repository replaced with ASCII punctuation.

### Fixed

- `/init` in OCI layers is matched correctly on Windows.
- App builds use a Linux cross toolchain instead of the host `llvm-gcc`;
  `doctor` reports a cross GCC found on PATH; `arch <name>` no longer
  exits 1 when crosstool-ng is absent.
- Kernel builds use Homebrew GNU make instead of macOS make 3.81.
- `init` creates the image's parent directory; runtime errors no longer
  print the whole usage text.
- `init` reopens a configured workspace image instead of creating a new one.
- Root help shows `elmos kernel build` instead of a command that does not
  exist.
- `assets/libraries/asm` symlinks to one machine's paths replaced with the
  headers themselves.
- CLI workflows aligned and the module manifest restored.
- Broken API link in the developer documentation index.

## [3.2.0] - 2026-01-20

See the [GitHub release](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.2.0).

## [3.1.0] - 2026-01-09

See the [GitHub release](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.1.0).

## [3.0.0] - 2026-01-01

See the [GitHub release](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.0.0).

## [2.1.0] - 2025-12-23

See the [GitHub release](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v2.1.0).

## [2.0.0] - 2025-12-20

See the [GitHub release](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v2.0.0).

## [1.1.0] - 2025-12-18

See the [GitHub release](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v1.1.0).

## [1.0.0] - 2025-12-14

See the [GitHub release](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v1.0.0).

[Unreleased]: https://github.com/NguyenTrongPhuc552003/elmos/compare/v3.2.0...HEAD
[3.2.0]: https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.2.0
[3.1.0]: https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.1.0
[3.0.0]: https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.0.0
[2.1.0]: https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v2.1.0
[2.0.0]: https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v2.0.0
[1.1.0]: https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v1.1.0
[1.0.0]: https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v1.0.0
