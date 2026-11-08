# Security Policy

ELMOS is research software maintained as part of a Master's thesis. It runs
build tools, mounts disk images and downloads release assets on the host, so
some bugs can have security impact. Do not report security-sensitive issues
through public issue trackers if they could expose users or downstream
systems.

For responsible disclosure, contact the project maintainer privately and
include:

- affected component or command
- elmos version (`elmos version`), macOS version and chip
- impacted workflow (for example kernel build, rootfs creation, demo
  download, QEMU run)
- reproduction steps
- expected security impact
- logs or traces needed to validate the issue
