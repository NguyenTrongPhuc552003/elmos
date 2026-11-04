package rootfs

import (
	"archive/tar"
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/NguyenTrongPhuc552003/elmos/third_party/ext4/compactext4"
)

// tarToExt4 writes the files of a flattened tar stream (no whiteouts left)
// into a new ext4 file system of `size` bytes on w, writable and with the
// space beyond the content free. Ownership, modes, device numbers,
// symbolic and hard links and extended attributes (PAX SCHILY.xattr.*)
// come from the tar headers, so no root is needed on the host.
//
// Adapted from hcsshim's tar2ext4.ConvertTarToExt4 (MIT, see
// third_party/ext4/LICENSE), without its VHD, dm-verity and whiteout
// options and their dependencies.
func tarToExt4(r io.Reader, w io.ReadWriteSeeker, size int64) error {
	t := tar.NewReader(bufio.NewReader(r))
	fs := compactext4.NewWriter(w, compactext4.DiskSize(size), compactext4.Writable)
	for {
		hdr, err := t.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		name := hdr.Name
		if err := fs.MakeParents(name); err != nil {
			return fmt.Errorf("parent directories of %s: %w", name, err)
		}
		if hdr.Typeflag == tar.TypeLink {
			if err := fs.Link(hdr.Linkname, name); err != nil {
				return fmt.Errorf("hard link %s: %w", name, err)
			}
			continue
		}
		f := &compactext4.File{
			Mode:     uint16(hdr.Mode),
			Atime:    hdr.AccessTime,
			Mtime:    hdr.ModTime,
			Ctime:    hdr.ChangeTime,
			Crtime:   hdr.ModTime,
			Size:     hdr.Size,
			Uid:      uint32(hdr.Uid),
			Gid:      uint32(hdr.Gid),
			Linkname: hdr.Linkname,
			Devmajor: uint32(hdr.Devmajor),
			Devminor: uint32(hdr.Devminor),
			Xattrs:   map[string][]byte{},
		}
		for key, value := range hdr.PAXRecords {
			const prefix = "SCHILY.xattr."
			if strings.HasPrefix(key, prefix) {
				f.Xattrs[key[len(prefix):]] = []byte(value)
			}
		}
		var typ uint16
		switch hdr.Typeflag {
		case tar.TypeReg:
			typ = compactext4.S_IFREG
		case tar.TypeSymlink:
			typ = compactext4.S_IFLNK
		case tar.TypeChar:
			typ = compactext4.S_IFCHR
		case tar.TypeBlock:
			typ = compactext4.S_IFBLK
		case tar.TypeDir:
			typ = compactext4.S_IFDIR
		case tar.TypeFifo:
			typ = compactext4.S_IFIFO
		default:
			continue // e.g. PAX global headers
		}
		f.Mode = f.Mode&^compactext4.TypeMask | typ
		if err := fs.Create(name, f); err != nil {
			return fmt.Errorf("create %s: %w", name, err)
		}
		if _, err := io.Copy(fs, t); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return fs.Close()
}
