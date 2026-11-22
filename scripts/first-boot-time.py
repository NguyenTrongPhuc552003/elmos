#!/usr/bin/env python3
"""Time ELMOS's first boot on this host: `elmos demo linux` until the guest
prints "System ready." (docs/developer/first-boot.md).

Usage: first-boot-time.py ELMOS_BINARY CACHE_DIR LABEL [TIMEOUT_S]
Prints one JSON line {label, ok, seconds_to_boot_start, seconds_to_ready,
reason} and appends a row to $GITHUB_STEP_SUMMARY when it is set. The demo
is stopped (with QEMU) once the marker appears. Standard library only;
runs on Linux, macOS and Windows.
"""
import json
import os
import signal
import subprocess
import sys
import threading
import time

BOOTING = "Booting (exit with Ctrl-A X)"
READY = "System ready."


def kill_tree(p):
    if os.name == "nt":
        subprocess.run(["taskkill", "/T", "/F", "/PID", str(p.pid)], capture_output=True)
    else:
        try:
            os.killpg(p.pid, signal.SIGKILL)
        except ProcessLookupError:
            pass


def main(argv):
    if len(argv) not in (3, 4):
        sys.exit(__doc__)
    binary, cache, label = argv[0], argv[1], argv[2]
    limit = float(argv[3]) if len(argv) == 4 else 900
    kw = {"creationflags": subprocess.CREATE_NEW_PROCESS_GROUP} if os.name == "nt" else {"start_new_session": True}
    start = time.monotonic()
    p = subprocess.Popen([binary, "demo", "linux", "--cache", cache], stdin=subprocess.DEVNULL,
                         stdout=subprocess.PIPE, stderr=subprocess.STDOUT, **kw)
    seen, tail = {}, []

    def read():
        for raw in iter(p.stdout.readline, b""):
            line = raw.decode(errors="replace").rstrip()
            tail.append(line)
            del tail[:-20]
            for key, marker in (("boot", BOOTING), ("ready", READY)):
                if marker in line and key not in seen:
                    seen[key] = round(time.monotonic() - start, 1)
            if "ready" in seen:
                return

    t = threading.Thread(target=read, daemon=True)
    t.start()
    t.join(limit)
    exited = p.poll()
    kill_tree(p)
    ok = "ready" in seen
    reason = "" if ok else (f"exited {exited} before the marker" if exited is not None else f"no marker in {limit:.0f} s")
    res = {"label": label, "ok": ok, "seconds_to_boot_start": seen.get("boot"),
           "seconds_to_ready": seen.get("ready"), "reason": reason, "last_lines": [] if ok else tail[-8:]}
    print(json.dumps(res))
    summary = os.environ.get("GITHUB_STEP_SUMMARY")
    if summary:
        with open(summary, "a") as f:
            f.write(f"| {label} | {'yes' if ok else 'no'} | {res['seconds_to_boot_start']} | "
                    f"{res['seconds_to_ready']} | {reason} |\n")
    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
