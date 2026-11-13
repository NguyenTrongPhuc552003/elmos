# Roadmap

Where ELMOS is going, and where it is not. Reviewed in March 2027 against
what people actually ask for.

## Who ELMOS is for

People meeting Embedded Linux for the first time (students, software
developers new to embedded work, teachers), and the HAVEN ecosystem
([HAVEN](https://github.com/NguyenTrongPhuc552003/haven),
[KDAL](https://github.com/NguyenTrongPhuc552003/kdal)), which needs a
reproducible kernel, root file system and QEMU setup.

## The niche: learn and prototype

The fastest path from nothing to a running embedded Linux, on the laptop
the learner already has (macOS, Windows, Linux), with each stage visible:
toolchain, kernel, root file system, boot in QEMU, a first module.

ELMOS is not a general distribution builder. If you need a product image
with many packages, licence reports or long-term maintenance, Buildroot
and the Yocto Project are the established tools for that.

## Now

- Keep the documentation current with what the code does.
- The SDK (kernel and root file system built by CI for KDAL and HAVEN)
  stays reproducible and its smoke boot ends on its own.
- A place for questions and first contributions: GitHub Discussions and
  the `good first issue` and `help wanted` labels.
- One documentation site (the project site, haven-docs) instead of two.

The SDK role for KDAL and HAVEN has priority over the rest.

## Later (not before the March 2027 review)

- A package catalogue beyond Alpine.
- Real boards beyond the CIX P1 board HAVEN targets.
- SBOM and licence reports.
- More kernel versions in the SDK.
- A path from an ELMOS project to a Buildroot defconfig or a Yocto layer.

## Measured, not claimed

How long the first boot takes on each host is to be measured before it
is stated anywhere.

## Tell us

Questions and ideas: [Discussions](https://github.com/NguyenTrongPhuc552003/elmos/discussions).
Bugs and concrete requests: [issues](https://github.com/NguyenTrongPhuc552003/elmos/issues).
