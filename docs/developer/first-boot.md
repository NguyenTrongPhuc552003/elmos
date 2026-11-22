# Time to first boot on three hosts

**Status:** design and expectations, 2026-11-22, before the first CI run
(owner-approved: the positioning's evidence item, "time-to-first-boot
measured on the three host systems").

ELMOS's positioning claims a first embedded Linux "in minutes with
prebuilt pieces" on the laptop the learner already has. This measures it
on GitHub's runners for Linux, macOS and Windows.

**What is measured.** `scripts/first-boot-time.py` runs `elmos demo
linux` with an empty cache (cold: download, unpack, boot) and again with
the cache filled (warm: boot only), and records the seconds until the
demo prints `Booting` (assets ready) and until the guest prints
`System ready.` (ELMOS's `/init`). It then stops the demo and QEMU.
Workflow: `.github/workflows/first-boot.yml`, on demand and weekly;
QEMU from each host's package manager (apt, Homebrew, Chocolatey).

**On the Mac (owner's, 2026-11-22, one run each)** [F]: cold 10.2 s to
`System ready.` (9.1 s of it fetching and unpacking), warm 0.9 s.

| Row | Host | Expected |
| - | - | - |
| F1 | Ubuntu 24.04 runner | cold and warm both reach `System ready.`; cold under 5 minutes |
| F2 | macOS runner (arm64, Homebrew QEMU) | the same |
| F3 | Windows runner (Chocolatey QEMU) | **fails before `System ready.`** [I]: the demo shares the modules over virtio-9p (`-fsdev local`), which QEMU's Windows builds do not include; the failure is the measurement, a Multi-OS gap to fix in `demo.LinuxArgs` |

Runner numbers are not laptop numbers (shared virtual machines, network
to GitHub's release storage); they show that the path works and its
order of magnitude.


## First run (2026-11-22, run 44994413431)

| Row | Host, QEMU | Cold: to `Booting` / to `System ready.` | Warm | Result |
| - | - | - | - | - |
| F1 | Ubuntu 24.04, QEMU 8.2.2 | 2.6 s / 3.9 s | 0.2 s / 1.5 s | holds |
| F2 | macOS (arm64), QEMU 11.1.3 | 2.5 s / 3.6 s | 0.2 s / 1.2 s | holds |
| F3 | Windows, QEMU 11.1.0 | 1.9 s / never | 0.1 s / never | holds as predicted: `There is no option group 'fsdev'` |

[F] from the job logs. Two things the run showed besides the rows:

- On Windows the demo did not end when QEMU refused the command line:
  the timing script waited its full limit (900 s, then 600 s) instead of
  stopping at the failure. Not examined yet [?].
- The fix for F3 (same pull request as these results): on Windows the
  demo leaves out the 9p share and says so; the guest boots without
  KDAL's module there. The next weekly run checks F3 again, expected to
  reach `System ready.`.
