# Changelog

User-facing changes in ELMOS. The full list is the repository's
[CHANGELOG.md](https://github.com/NguyenTrongPhuc552003/elmos/blob/main/CHANGELOG.md);
each release's notes are on its
[GitHub release page](https://github.com/NguyenTrongPhuc552003/elmos/releases).

## Unreleased (since v3.2.0)

- `elmos demo linux` and `elmos demo haven`: prebuilt demos that need only
  QEMU (no kernel build, no workspace volume, no sudo); assets are checked
  against the release's SHA256SUMS and cached for offline use.
- The OCI root file system engine is the default for `rootfs create`: it
  writes ext4 in pure Go without extracting onto the Mac; debootstrap
  stays available with `--engine debootstrap`, and e2fsprogs is optional.
- Kernel patches and the Linux headers macOS lacks are embedded in the
  binary; `kernel clone` (shallow v6.18 by default) applies the patches.
- Taskfile tasks to build, test and boot HAVEN (`haven:*`).
- CI on Linux, macOS and Windows, and an SDK (Linux kernel, root file
  system) built by CI and published privately for KDAL's and HAVEN's CI
  ([SDK built by CI](../developer/sdk-ci.md)).
- Fixes: `/init` in OCI layers on Windows, app builds with a Linux cross
  toolchain, kernel builds with GNU make, `init` reopening a workspace.

## Releases

| Version | Date |
| ------- | ---- |
| [3.2.0](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.2.0) | 2026-01-20 |
| [3.1.0](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.1.0) | 2026-01-09 |
| [3.0.0](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v3.0.0) | 2026-01-01 |
| [2.1.0](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v2.1.0) | 2025-12-23 |
| [2.0.0](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v2.0.0) | 2025-12-20 |
| [1.1.0](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v1.1.0) | 2025-12-18 |
| [1.0.0](https://github.com/NguyenTrongPhuc552003/elmos/releases/tag/v1.0.0) | 2025-12-14 |
