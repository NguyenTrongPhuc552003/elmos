# Linux demo release built by CI (ecosystem step E7)

**Status:** owner-approved plan (2026-11-13). Design and expectations,
written and merged before the workflow exists; results are added below by
a later pull request, expectations are never edited after them. Labels:
**[F]** fact from code, records, CI logs or release metadata, **[O]**
observation made for this record, **[I]** inference, **[P]** proposal.
The HAVEN half (the try bundle) has its own record, haven
`docs/methodology/stage2/20261113-e7-releases/README.md`.

## Why

- [F] `elmos demo linux` downloads the release `demo-linux-v6.18-1`
  (`core/domain/demo/demo.go`, `LinuxTag`). It was built by hand on the
  owner's Mac on 2026-11-04 (its `MANIFEST.txt`): Linux v6.18 with
  Homebrew LLVM, release string `6.18.0-dirty`, `kdal.ko` from KDAL
  `6701e5a`, the device tree dumped from the Mac's QEMU.
- [F] The kernel banner in its `Image.gz` carries the Mac's login and host
  name (checked 2026-11-13, haven `docs/methodology/LOG-SANITIZING.md`).
  The names are not written here.
- [F] The same kernel is now built by CI with neutral names (`haven@ci`)
  and published as the private package `elmos-sdk`
  (`docs/developer/sdk-ci.md`); HAVEN's CI pins one build of it, and the
  KDAL commit, in haven's `ecosystem.lock.json`.

## What the release holds

The same six file names as `demo-linux-v6.18-1`, so `elmos demo linux`
needs no change except the tag:

| File | Made from |
| - | - |
| `Image.gz` | the SDK package's `Image.gz`, unchanged |
| `rootfs.ext4.gz` | the SDK package's `rootfs.ext4.gz` (alpine:3, 1 GiB, ELMOS's `/init`), unchanged |
| `modules.tar.gz` | `kdal/kdal.ko` built at the pinned KDAL commit against the SDK kernel (KDAL's `scripts/ci/kernel-prepare.sh`, as HAVEN's CI does), debug information stripped; `guesync.sh` from `sdk/demo/guesync.sh` (the script of `demo-linux-v6.18-1`, now under version control unchanged); entries owned by uid 0 and gid 0 with no names |
| `virt-kdal.dtb` | QEMU's own device tree for the demo machine (`virt`, cortex-a72, 2 CPUs, 1 GiB, dumped with `dumpdtb` on the runner), the random `kaslr-seed` and `rng-seed` removed from `/chosen` so the file is the same on every run, and a node `kdal-gpio@9030000` (`compatible = "kdal,device"`, the PL061's `reg` and `interrupts`) added with `fdtput`, as in the old file |
| `MANIFEST.txt` | what each file is, the haven commit the pins came from, the SDK tag and digest, the KDAL commit, the QEMU version, the run |
| `SHA256SUMS` | the five files above |

## Workflow

`.github/workflows/demo-release.yml`, `workflow_dispatch` only, inputs
`tag` (form `demo-linux-v6.18-N`) and `haven_ref` (default `main`):

1. **Pins** from haven's `ecosystem.lock.json` at `haven_ref`, resolved to
   a commit first (the file is read at that commit, which goes into
   `MANIFEST.txt`): SDK package, tag, digest, ELMOS commit and run; KDAL
   repository and commit; QEMU version. The SDK is pulled by digest; the
   tag must resolve to the same digest; the package's `manifest.json` must
   name the pinned ELMOS commit and run; the runner's QEMU must be the
   pinned version. Any disagreement fails the run.
2. **Build** (`ubuntu-24.04`, `contents: write`, `packages: read`):
   `kdal.ko`, the device tree, the archive, the manifest and the hashes.
   `modinfo -F vermagic kdal.ko` must start with the SDK's
   `kernel.release`.
3. **Boot check** on the runner: `go run ./sdk/demoboot DIR` prepares the
   directory with the same code as `elmos demo linux` (the new
   `demo.PrepareLinux`, moved out of the command unchanged) and runs QEMU
   with `demo.LinuxArgs`, reading the console until it has seen every
   marker: `Linux version 6.18.0 (haven@ci)`, `Run /init as init process`,
   `[GUEST] kdal loaded OK`, `KDAL resource demand`, `System ready.`; then
   it stops QEMU. It fails if a marker is missing after 300 s. (The check
   does not use `elmos_poweroff`: it keeps the user's kernel command line
   of `LinuxArgs` unchanged and stops QEMU itself; a user gets the shell
   as before.)
4. **Name check:** `sdk/check-public-assets.py` (the same script as
   haven's `scripts/release/check-public-assets.py`, rules there) on the
   assets and the release notes: no tar entry owned by anyone but uid 0,
   no `/Users/` path, every `Linux version` banner built by `haven@ci`,
   and no word of `RELEASE_DENY` when that is set (it is set only on the
   Mac, at run time, from `id -un` and the host name; never written down).
5. **Draft:** a draft pre-release for `tag` at the dispatched commit,
   title ending in `[untested]`. A draft is visible only to people with
   write access and creates no Git tag before it is published.
6. **test-macos** (`macos-15`, Homebrew QEMU): downloads the draft's
   assets through the API with `GITHUB_TOKEN`, checks `SHA256SUMS` and
   runs the same boot check; this is the host the demo is for.
7. **mark:** edits the draft's title and notes with the results; the
   title loses `[untested]` only if both boot checks passed, otherwise it
   ends in `[FAILED: do not publish]`.

Nothing publishes, deletes or uploads an Actions artifact (artifacts of a
public repository can be downloaded by any signed-in user). Publication
needs the owner's explicit yes on the specific draft; the coordinator asks
for it with the tag, commit, hashes, results and a local name check of the
downloaded draft. Afterwards a pull request moves `LinuxTag` (and, for the
HAVEN half, `HavenTag` and `HavenBundle`) to the new releases; it is
prepared after the drafts exist and merged only after publication. The
owner then deletes `demo-linux-v6.18-1` on GitHub and Codeberg.

[F] The KDAL module's `modinfo` author field is KDAL's author, as written
in KDAL's public source (`MODULE_AUTHOR`). It is declared authorship, not
a build-machine name, and the check does not flag it. [P] The owner
decides whether that is acceptable in a public binary (it was already in
`demo-linux-v6.18-1`).

## Expectations (before the first run)

| Row | Check | Expected |
| - | - | - |
| D1 | pins step | reads SDK `linux-6.18-5`, digest `sha256:cea68ea3...`, ELMOS `8ae1437` run 5, KDAL `f11e87e`, QEMU 8.2.2 (haven `main` after haven #89); the tag resolves to the digest; the manifest and QEMU checks pass |
| D2 | `kdal.ko` | builds; `vermagic` starts with `6.18.0 ` (the SDK's release, no `-dirty`) |
| D3 | `virt-kdal.dtb` | has `/kdal-gpio@9030000` with `compatible = "kdal,device"`, `reg = <0 0x9030000 0 0x1000>` and `interrupts = <0 7 4>` as in the old file; `/chosen` has no `kaslr-seed` and no `rng-seed`. Other differences from the old file (made by a newer QEMU on the Mac) are recorded, not predicted |
| D4 | boot check on the runner (QEMU 8.2.2) | all five markers, in under 120 s (the time is a prediction; 300 s is the limit) |
| D5 | name check in CI | exit 0; notes, if any, only `/home/runner` and `/home/buildozer` (Alpine's package builder, inside the root filesystem) |
| D6 | draft | `draft: true`, `prerelease: true`, the six files, target the dispatched commit; no Git tag of that name after the run; nothing published; no Actions artifact |
| D7 | test-macos (Homebrew QEMU, newer than 8.2.2) | all five markers: the device tree dumped by QEMU 8.2.2 works with the newer QEMU (prediction, the main risk) |
| D8 | mark | title without `[untested]`; notes give both QEMU versions and boot times |
| D9 | local name check on the Mac | exit 0 on the draft's files; exit 1 on `demo-linux-v6.18-1`, at rule (c) (its banner) and rule (d) |
| D10 | `go test ./...` with `demo.PrepareLinux` | passes; `LinuxArgs` unchanged |

Not shown here: that `elmos demo linux` fetches the published release
(the tag pull request checks that, after publication), Windows hosts.
