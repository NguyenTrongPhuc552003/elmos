# SDK built by CI (ecosystem step E2)

**Status:** design and expectations, 2026-11-09, written before the
workflow exists (owner-approved step E2 of the HAVEN ecosystem plan; owner
rules of 2026-11-09: CI first, CI outputs private, protected `main`).

## Why

Today the kernel that KDAL's CI and HAVEN's Linux runs use was built by
ELMOS on the owner's Mac and uploaded by hand as public releases
(`demo-linux-v6.18-1` here, `ci-kernel-v6.18-kunit` in KDAL). The Mac's
resources are limited, a hand-made release cannot be rebuilt by anyone
else, and its version string carries the build machine's user and host
name. The SDK moves that build to GitHub Actions and publishes it where
only the three project repositories can read it.

## What the workflow does

`.github/workflows/sdk.yml`, on `workflow_dispatch` and on pushes to `main`
that change `sdk/` or the workflow (not on every push: a kernel build is
long):

1. **Kernel:** Linux 6.18 from kernel.org (tarball hash checked against the
   value pinned in `sdk/versions.env`), arm64 `defconfig` plus the fragment
   `sdk/kernel.cfg` (`CONFIG_KUNIT=y`, as today's releases), built with
   LLVM natively on GitHub's arm64 runner (`ubuntu-24.04-arm`), with
   `KBUILD_BUILD_USER=haven`, `KBUILD_BUILD_HOST=ci` and a fixed
   `KBUILD_BUILD_TIMESTAMP`. Outputs: `Image.gz`, `config`,
   `Module.symvers`, `kernel.release` (the files KDAL's
   `scripts/ci/kernel-prepare.sh` reads today).
2. **Root filesystem:** `elmos rootfs create` (OCI engine, `alpine:3`, the
   same `/init`), built from this commit, in a fresh ELMOS workspace on the
   runner; output `rootfs.ext4.gz`.
3. **Smoke boot:** the kernel and root filesystem in QEMU `virt`
   (cortex-a72, two CPUs, TCG) on the runner, until the kernel runs `/init`.
4. **Manifest:** `manifest.json` (Linux version and tarball hash, config
   hash, ELMOS commit, build user and host, file hashes) and `SHA256SUMS`.
5. **Publish:** the files as one OCI artifact with `oras` to
   `ghcr.io/nguyentrongphuc552003/elmos-sdk:linux-6.18-<run number>` (and
   `:linux-6.18`), with `GITHUB_TOKEN` (`packages: write`). A GHCR package
   is private when first published and does not take the repository's
   visibility; it is linked to this repository by the
   `org.opencontainers.image.source` annotation so it inherits its access
   permissions. KDAL and HAVEN read it after the owner grants them Read
   access in the package settings (step E3 uses that).
6. **Privacy check:** an anonymous fetch of the manifest is expected to be
   refused, and the workflow fails if it is not.

## Expectations (before the first run)

| Check | Expected |
| - | - |
| S1: dispatched run on `ubuntu-24.04-arm` | succeeds; the kernel build takes under 90 minutes (prediction, recorded either way) |
| S2: version string in the built `Image` | `Linux version 6.18.0` with `(haven@ci)`; no personal user or host name |
| S3: root filesystem on a Linux runner | `elmos rootfs create` works there (risk: ELMOS was used on macOS until now; if it fails, the reason is recorded and the rootfs step is fixed or deferred, not faked) |
| S4: smoke boot | the kernel log shows `Run /init as init process` within 5 minutes |
| S5: package | `elmos-sdk` exists in GHCR with visibility `private` (GitHub API) and the run's tag |
| S6: anonymous fetch | refused (HTTP 401 or 403) |
| S7: contents | `SHA256SUMS` matches the files; `manifest.json` names this commit |
| S8: nothing public | no release, no Actions artifact, no public link is created |

Results are added below after the first run.
