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
