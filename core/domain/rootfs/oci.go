package rootfs

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Microsoft/hcsshim/ext4/tar2ext4"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"

	"codeberg.org/NguyenTrongPhuc552003/elmos/assets"
)

// The OCI engine builds the root filesystem from a container image instead
// of debootstrap. The image is pulled for the target architecture, its
// layers are flattened into one tar stream, and that stream is written
// straight into an ext4 image. Nothing is extracted onto the macOS file
// system, so file ownership, modes and case-sensitive names survive without
// root, fakeroot or a case-sensitive volume.

// ociPlatform maps an elmos architecture to an OCI platform.
func ociPlatform(arch string) (*v1.Platform, error) {
	switch arch {
	case "arm64", "":
		return &v1.Platform{OS: "linux", Architecture: "arm64"}, nil
	case "arm":
		return &v1.Platform{OS: "linux", Architecture: "arm", Variant: "v7"}, nil
	case "riscv":
		return &v1.Platform{OS: "linux", Architecture: "riscv64"}, nil
	default:
		return nil, fmt.Errorf("no OCI platform for architecture %q", arch)
	}
}

// pullFlattened pulls ref for the platform and returns its flattened file
// system as one tar stream.
func pullFlattened(ctx context.Context, ref string, platform *v1.Platform) (io.ReadCloser, error) {
	img, err := crane.Pull(ref, crane.WithPlatform(platform), crane.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("pull %s (%s/%s): %w", ref, platform.OS, platform.Architecture, err)
	}
	return mutate.Extract(img), nil
}

// withInit copies the image tar to w, replacing any /init with the elmos
// one, and closes the archive.
func withInit(w io.Writer, image io.Reader, initScript []byte) error {
	tr := tar.NewReader(image)
	tw := tar.NewWriter(w)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if strings.TrimPrefix(filepath.Clean("/"+hdr.Name), "/") == "init" {
			continue
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := io.Copy(tw, tr); err != nil {
			return err
		}
	}
	if err := tw.WriteHeader(&tar.Header{
		Name:     "init",
		Mode:     0o755,
		Size:     int64(len(initScript)),
		Typeflag: tar.TypeReg,
		ModTime:  time.Now(),
	}); err != nil {
		return err
	}
	if _, err := tw.Write(initScript); err != nil {
		return err
	}
	return tw.Close()
}

// parseSize turns "5G", "512M" or a byte count into bytes.
func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	mult := int64(1)
	switch {
	case strings.HasSuffix(s, "G"):
		mult, s = 1<<30, strings.TrimSuffix(s, "G")
	case strings.HasSuffix(s, "M"):
		mult, s = 1<<20, strings.TrimSuffix(s, "M")
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid size %q", s)
	}
	return n * mult, nil
}

// e2fsTool returns the path of an e2fsprogs tool. Homebrew installs
// e2fsprogs keg-only, so its sbin is not on PATH by default.
func (c *Creator) e2fsTool(ctx context.Context, name string) string {
	out, err := c.exec.Output(ctx, "brew", "--prefix", "e2fsprogs")
	if err == nil {
		if prefix := strings.TrimSpace(string(out)); prefix != "" {
			return filepath.Join(prefix, "sbin", name)
		}
	}
	return name
}

// createFromOCI writes diskImage from an OCI image.
func (c *Creator) createFromOCI(ctx context.Context, diskImage, ref, size string) error {
	if ref == "" {
		ref = DefaultOCIImage
	}
	sizeBytes, err := parseSize(size)
	if err != nil {
		return err
	}
	platform, err := ociPlatform(c.cfg.Build.Arch)
	if err != nil {
		return err
	}
	initScript, err := assets.GetInitScript()
	if err != nil {
		return err
	}
	pull := c.pull
	if pull == nil {
		pull = pullFlattened
	}
	image, err := pull(ctx, ref, platform)
	if err != nil {
		return err
	}
	defer image.Close()

	if err := os.MkdirAll(filepath.Dir(diskImage), 0o755); err != nil {
		return err
	}
	tmp := diskImage + ".partial"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer os.Remove(tmp)

	pr, pw := io.Pipe()
	go func() { pw.CloseWithError(withInit(pw, image, initScript)) }()
	convErr := tar2ext4.Convert(pr, f, tar2ext4.MaximumDiskSize(sizeBytes))
	_ = pr.CloseWithError(convErr)
	if cerr := f.Close(); convErr == nil {
		convErr = cerr
	}
	if convErr != nil {
		return fmt.Errorf("write ext4 from %s: %w", ref, convErr)
	}

	// tar2ext4 writes a compact, read-only file system with no free space
	// and no journal. Make it an ordinary writable ext4 of the asked size.
	if err := c.fsck(ctx, tmp); err != nil {
		return err
	}
	steps := [][]string{
		{c.e2fsTool(ctx, "tune2fs"), "-O", "^read-only", tmp},
		{c.e2fsTool(ctx, "resize2fs"), tmp, size},
		{c.e2fsTool(ctx, "tune2fs"), "-j", "-L", "rootfs", tmp},
	}
	for _, s := range steps {
		if err := c.exec.Run(ctx, s[0], s[1:]...); err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(s[0]), err)
		}
	}
	if err := c.fsck(ctx, tmp); err != nil {
		return err
	}
	return os.Rename(tmp, diskImage)
}

// fsck runs e2fsck -fy. Exit status 1 means errors were found and
// corrected, which is expected on the compact image (padding bitmap).
func (c *Creator) fsck(ctx context.Context, image string) error {
	err := c.exec.Run(ctx, c.e2fsTool(ctx, "e2fsck"), "-fy", image)
	var exitErr *exec.ExitError
	if err == nil || (errors.As(err, &exitErr) && exitErr.ExitCode() == 1) {
		return nil
	}
	return fmt.Errorf("e2fsck: %w", err)
}
