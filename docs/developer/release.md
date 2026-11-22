# Releases built by CI

**Status:** design and expectations, 2026-11-22, before the first run.
The owner said yes to v3.3.0 (2026-11-22); the earlier releases (v1.0.0
to v3.2.0) were built by hand.

`.github/workflows/release.yml`, `workflow_dispatch` with the input
`version` (for example `3.3.0`):

1. **build** (Ubuntu): the tag `vX.Y.Z` must not exist and `CHANGELOG.md`
   must have a section `## [X.Y.Z] - DATE`. `elmos` is cross-compiled with
   `CGO_ENABLED=0 -trimpath` for darwin, linux and windows on arm64 and
   amd64, with `Version`, `Commit` and `BuildDate` set; `SHA256SUMS`
   covers the six binaries. The notes are the changelog section. The
   binaries and notes pass `sdk/check-public-assets.py` (no `/Users/`
   path, no private word). A **draft** release is created for the
   dispatched commit.
2. **smoke**: on Ubuntu, macOS and Windows runners, the host's binary is
   downloaded from the draft, its hash checked, `elmos version` must name
   the version, and `elmos doctor` runs (its report is not the check).
3. **mark**: the draft's title loses `[untested]` only if every smoke job
   passed, otherwise it ends in `[FAILED: do not publish]`.

Publishing the draft is a separate step on the owner's yes. Afterwards the
Homebrew formula moves to the new tag and its hash.

| Row | Check | Expected |
| - | - | - |
| R1 | build | six binaries, `SHA256SUMS`, the name check passes, a draft for `v3.3.0` at the dispatched commit, no Git tag yet |
| R2 | smoke on the three runners | each hash matches; `elmos version` prints `v3.3.0` |
| R3 | mark | title `v3.3.0: Embedded Linux on Multi-OS` |

**Order for v3.3.0** [P]: the replacement Linux demo release
(`demo-release.md`) is published first and `demo.LinuxTag` moved to it,
so v3.3.0 does not ship a binary that downloads the old demo, whose
kernel banner carries a private name.

## v3.3.0 (2026-11-22, run 45041327081 at `171f750`) [F]

| Row | Result |
| - | - |
| R1 | six binaries and `SHA256SUMS`; the name check `0 failure(s), 0 note(s)`; a draft `v3.3.0` at `171f750`; no Git tag before publication. **Holds** |
| R2 | each downloaded binary's hash matches on ubuntu-24.04, macos-15 and windows-latest; `elmos version` prints `Version: v3.3.0`, `Commit: 171f750...`. **Holds** |
| R3 | title `v3.3.0: Embedded Linux on Multi-OS`. **Holds** |

Published on the owner's yes of 2026-11-22 (16:04 UTC, marked latest);
the tag `v3.3.0` points at `171f750`, after `demo.LinuxTag` moved to
`demo-linux-v6.18-2`.

The Homebrew formula, found while moving it to v3.3.0: `go build` with
`std_go_args` builds the repository root, which has had no Go files since
the main package moved to `cmd/elmos`, and the version came from `git
describe`, which a release tarball cannot run. Both fixed; the build
command was reproduced from the v3.3.0 tarball in the workspace (`elmos
version` prints `v3.3.0`), and the formula passes `ruby -c`. Not run:
`brew install` (it writes outside the workspace; the owner's step). The
formula's comment named a tap repository that does not exist; it now
gives the tap's URL. The two `release:darwin` targets in `Taskfile.yml`
built the root the same way and are fixed too.
