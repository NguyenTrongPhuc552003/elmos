#!/usr/bin/env bash
# Check that ELMOS's kernel patches apply to the kernel version they are
# for, and that the host tool they patch still builds and works
# (audit item E17). Fetches only the files the patches touch, from the
# kernel.org git tree at the version's tag (GitHub's raw host rate-limits).
# Usage: scripts/check-kernel-patches.sh [VERSION]   (default v6.18)
set -euo pipefail
version="${1:-v6.18}"
root="$(cd "$(dirname "$0")/.." && pwd)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

found=0
for p in "$root"/assets/patches/"$version"/*/*.patch; do
	[ -e "$p" ] || continue
	found=1
	sed -n 's|^+++ b/||p' "$p" | while read -r f; do
		mkdir -p "$work/linux/$(dirname "$f")"
		curl -fsSL --retry 3 -o "$work/linux/$f" \
			"https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/plain/$f?h=$version"
	done
	echo "applying $(basename "$p")"
	patch -d "$work/linux" -p1 --forward -i "$p"
done
[ "$found" = 1 ] || { echo "no patches for $version"; exit 1; }

cc -Wall -O2 -o "$work/gen_init_cpio" "$work/linux/usr/gen_init_cpio.c"
printf 'hello\n' > "$work/hello"
printf 'dir /dev 0755 0 0\nfile /hello %s 0644 0 0\n' "$work/hello" > "$work/list"
"$work/gen_init_cpio" "$work/list" > "$work/out.cpio"
(cd "$work" && mkdir x && cd x && cpio -id < ../out.cpio 2>/dev/null)
cmp "$work/hello" "$work/x/hello"
echo "patches apply to $version; gen_init_cpio builds and copies file data"
