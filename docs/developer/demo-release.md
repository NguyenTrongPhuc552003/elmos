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

## Corrections before the first run (2026-11-22)

The workflow is written nine days after this design (it had not been,
though the design was merged on 2026-11-13; found while timing the demo,
`first-boot.md`). Before its first run:

- **D1**: haven's `ecosystem.lock.json` now pins KDAL `5622c47`, not
  `f11e87e`; the SDK pins (`linux-6.18-5`, `sha256:cea68ea3...`, ELMOS
  `8ae1437`, run 5) and QEMU 8.2.2 are unchanged. D1 reads the lock at
  `haven_ref`, so it expects `5622c47` now.
- The boot check is `sdk/demoboot` (`go run ./sdk/demoboot DIR`), the
  name check `sdk/check-public-assets.py` (a copy of haven's), the guest
  script `sdk/demo/guesync.sh` (the old release's, unchanged), and
  `demo.PrepareLinux` holds the unpacking `elmos demo linux` did (D10).
- Checked on the Mac before the workflow ran: `demoboot` on the old
  release sees every marker except the banner (`Linux version 6.18.0
  (haven@ci)`), and the name check fails the old release at rule (c)
  (its banner) and rule (a) (its tar entries carry owner names). D9's
  "rule (d)" needs `RELEASE_DENY`, which is set only for the local check.

## Run 1 (2026-11-22, run 44997203315): failed at D3

D1 and D2 passed (pins at haven `c48b62e`: SDK `linux-6.18-5`,
`sha256:cea68ea3...`, KDAL `5622c47`, QEMU 8.2.2). D3 stopped at the
device-tree dump: `qemu-system-aarch64: failed to find romfile
"efi-virtio.rom"`. Without `-nic`, QEMU adds its default network card,
a virtio-net PCI device whose boot ROM comes from a package that
`--no-install-recommends` leaves out. The dump now passes `-nic none`.
The PCI card is not in the device tree (only the PCIe host bridge is,
and it is there either way), and the demo itself uses
`virtio-net-device` (MMIO, no ROM), so D3's expectations are unchanged.
Nothing was drafted; the tag is free for run 2.

## Run 2 (2026-11-22, run 44997888799 at `3d55296`): every row holds

[F] from the run's log, the draft's API record and its files downloaded
on the Mac.

| Row | Result |
| - | - |
| D1 | pins from haven `main` as corrected: SDK `linux-6.18-5`, `sha256:cea68ea3...`, KDAL `5622c47`, QEMU 8.2.2; the tag resolves to the digest, the manifest and QEMU checks pass. **Holds** |
| D2 | `vermagic: 6.18.0 SMP preempt mod_unload aarch64`. **Holds** |
| D3 | `/kdal-gpio@9030000` with `compatible = "kdal,device"`, `reg = <0x00 0x9030000 0x00 0x1000>`, `interrupts = <0x00 0x07 0x04>`; `/chosen` holds only `stdout-path`. **Holds** |
| D4 | every marker in 1.3 s on the Linux runner (QEMU 8.2.2). **Holds**; the 120 s prediction was far too cautious |
| D5 | exit 0, one note: `/home/runner` in `kdal.ko`. **Holds** |
| D6 | `draft: true`, `prerelease: true`, the six files, target `3d55296`; no Git tag `demo-linux-v6.18-2`, no Actions artifact. **Holds** |
| D7 | every marker in 1.7 s on the macOS runner (Homebrew QEMU 11.1.3): QEMU 8.2.2's device tree works with the newer QEMU. **Holds** |
| D8 | title `ELMOS Linux demo demo-linux-v6.18-2`, no `[untested]`; notes give both QEMU versions and times. **Holds** |
| D9 | on the Mac, with `RELEASE_DENY` set only for the command: the draft's files pass `shasum -c` and the check exits 0 (the same one note); `demo-linux-v6.18-1` exits 1 at rules (c) (its banner) and (d) (both denied words), and also at (a) (tar entries with owner names) and (b) (a macOS home path in `kdal.ko`). **Holds**, with two more rules than predicted |
| D10 | `go test ./...` passes in CI on #40 and #42. **Holds** |

SHA-256 of the draft's files:

    9421fa520e57504180c1d5d90cd73883dc34ec37e6f86bd8fecb581c2296aced  Image.gz
    e1569d0c58132d586d8cc0116fb24d3774042d4133b56fc51b674780fe975ba6  rootfs.ext4.gz
    ce0b0493895c3616b9afda6a0caabb3447cf2725217d4856f0062fd2695e22ab  modules.tar.gz
    c7b4ac8db79d2175f4a031fc5571d6f9ccaac9d8b1355b5c28a1f6b4218dce2a  virt-kdal.dtb
    265a165fa9c091e6ae3aafdcee3f1ad31c950587ebf68717de9b3c6239353a07  MANIFEST.txt

Next, in this order: the owner's yes on this draft, publication (which
creates the tag), a pull request moving `demo.LinuxTag` to
`demo-linux-v6.18-2` (it checks that `elmos demo linux` fetches the
published release), then the owner deletes `demo-linux-v6.18-1` on GitHub
and Codeberg.

## Publication (2026-11-22)

The owner said yes to this draft ("Approve", 2026-11-22, 22:53 ICT).
Published at 15:53 UTC; the tag `demo-linux-v6.18-2` points at
`3d55296`. `demo.LinuxTag` moves to it in the same pull request as this
section. Checked on the Mac before that pull request [F]: an `elmos`
built with the new tag and an empty cache fetched the published release
from GitHub, `shasum -c` passed on all five files, and the guest reached
`System ready.` 7.5 s after the command started (download included).
Left for the owner: deleting `demo-linux-v6.18-1` on GitHub and Codeberg.

## Old release deleted (2026-11-22)

The owner deleted `demo-linux-v6.18-1` on GitHub (the API answers 404)
and all Codeberg repositories, so nothing is left to delete. The demo
fetcher no longer falls back to Codeberg.
