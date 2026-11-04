# ext4 writer (from hcsshim)

`compactext4/compact.go` and `format/format.go` are copied from
[hcsshim](https://github.com/microsoft/hcsshim) v0.15.0
(`ext4/internal/compactext4`, `ext4/internal/format`), MIT licence
(`LICENSE`, Copyright (c) 2015 Microsoft). The upstream packages are
internal, so they cannot be imported, and the public `tar2ext4` package
pulls in hcsshim's logging, protobuf and telemetry dependencies.

elmos changes, all in `compact.go` and marked "elmos":

- `DiskSize`: the file system is as large as asked, with the space beyond
  the content free, and one inode per 16 KiB (the mke2fs default) instead
  of only as many as the content needs. A file is extended with
  `Truncate`, so the free space stays sparse.
- `Writable`: no read-only compatibility feature; a random UUID and the
  label `rootfs`.
- Inode bitmap padding bits are set (e2fsck reported "Padding at end of
  inode bitmap is not set").
- The `internal/memory` import is replaced by a local constant.

The result has no journal; Linux mounts it read-write ("without
journal"). `core/domain/rootfs/tarext4.go` is the tar-to-ext4 loop,
adapted from `tar2ext4.ConvertTarToExt4`. The rootfs tests run
`e2fsck -fn` on a written image when e2fsck is installed.
