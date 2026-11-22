
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
