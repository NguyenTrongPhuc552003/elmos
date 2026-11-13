// Package demo runs the prebuilt quick-start demos: no kernel build, no
// workspace volume, no sudo. Assets come from GitHub releases, are
// checked against the release's SHA256SUMS and cached under ~/.elmos/demo.
package demo

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ReleaseBase is where the demo releases live.
	ReleaseBase = "https://github.com/NguyenTrongPhuc552003"
	// MirrorBase holds the same releases on Codeberg, the project's earlier
	// home; used when a release is not (yet) on GitHub.
	MirrorBase = "https://codeberg.org/NguyenTrongPhuc552003"
	// LinuxTag is the elmos release with the Linux demo: a v6.18 arm64
	// kernel, an alpine:3 root file system and the KDAL module built for
	// that kernel.
	LinuxTag = "demo-linux-v6.18-1"
	// HavenTag and HavenBundle name the HAVEN try bundle release.
	HavenTag    = "try-20261103"
	HavenBundle = "haven-try-v0.6.0-162-gab6f825.tar.gz"
)

// Fetcher downloads release assets into a cache directory.
type Fetcher struct {
	Client   *http.Client
	CacheDir string
	Base     string // ReleaseBase unless a test overrides it
	Mirrors  []string
}

// NewFetcher returns a Fetcher caching under dir.
func NewFetcher(dir string) *Fetcher {
	return &Fetcher{Client: http.DefaultClient, CacheDir: dir, Base: ReleaseBase, Mirrors: []string{MirrorBase}}
}

func (f *Fetcher) get(ctx context.Context, url, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := f.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: HTTP %d", url, resp.StatusCode)
	}
	tmp := dst + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, dst)
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// parseSums reads "HASH  NAME" lines (sha256sum format).
func parseSums(path string) (map[string]string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	sums, order := map[string]string{}, []string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 || len(fields[0]) != 64 {
			continue
		}
		name := strings.TrimPrefix(fields[1], "*")
		// a plain file name: Release joins it to its cache directory
		// (filepath.Base("/") is "/", so a Base comparison alone lets "/" in)
		if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
			return nil, nil, fmt.Errorf("SHA256SUMS: unsafe file name %q", name)
		}
		sums[name] = fields[0]
		order = append(order, name)
	}
	return sums, order, sc.Err()
}

// Release downloads every file listed in the release's SHA256SUMS (kept
// if already cached with the right hash) and returns the directory.
// Without network it uses the cache, still checked against the cached
// SHA256SUMS.
func (f *Fetcher) Release(ctx context.Context, repo, tag string) (string, error) {
	dir := filepath.Join(f.CacheDir, repo, tag)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	sumsPath := filepath.Join(dir, "SHA256SUMS")
	var base string
	var getErr error
	for _, b := range append([]string{f.Base}, f.Mirrors...) {
		base = fmt.Sprintf("%s/%s/releases/download/%s/", b, repo, tag)
		if getErr = f.get(ctx, base+"SHA256SUMS", sumsPath); getErr == nil {
			break
		}
	}
	if getErr != nil {
		if _, statErr := os.Stat(sumsPath); statErr != nil {
			return "", getErr
		}
	}
	sums, order, err := parseSums(sumsPath)
	if err != nil {
		return "", err
	}
	if len(order) == 0 {
		return "", fmt.Errorf("%s %s: SHA256SUMS lists no files", repo, tag)
	}
	for _, name := range order {
		path := filepath.Join(dir, name)
		if got, err := fileSHA256(path); err == nil && got == sums[name] {
			continue
		}
		if err := f.get(ctx, base+name, path); err != nil {
			return "", err
		}
		got, err := fileSHA256(path)
		if err != nil {
			return "", err
		}
		if got != sums[name] {
			os.Remove(path)
			return "", fmt.Errorf("%s: sha256 %s, release says %s", name, got, sums[name])
		}
	}
	return dir, nil
}

// Gunzip writes src (gzip) to dst unless dst is already newer than src.
func Gunzip(src, dst string) error {
	if si, err := os.Stat(src); err == nil {
		if di, err := os.Stat(dst); err == nil && !di.ModTime().Before(si.ModTime()) {
			return nil
		}
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	zr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	out, err := os.Create(dst + ".part")
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, zr); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Rename(dst+".part", dst)
}

// Untar extracts a .tar.gz into dir, refusing entries that leave dir.
func Untar(src, dir string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	zr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)
	root, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		dst := filepath.Join(root, h.Name)
		if dst != root && !strings.HasPrefix(dst, root+string(os.PathSeparator)) {
			return fmt.Errorf("archive entry %q leaves the target directory", h.Name)
		}
		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(dst, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(h.Mode)&0o755|0o600)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			if err := out.Close(); err != nil {
				return err
			}
		}
	}
}

// LinuxArgs is the QEMU command line for the Linux demo in dir: the
// kernel, the root file system with -snapshot (writes are thrown away, so
// the cached image stays clean) and the modules share that /init loads.
// If the release ships virt-kdal.dtb (QEMU's device tree for exactly this
// machine plus a "kdal,device" node over the PL061 GPIO), it is used, so
// KDAL finds a device-tree device and reports its resource demand.
func LinuxArgs(dir string) []string {
	args := []string{
		"-M", "virt", "-cpu", "cortex-a72", "-smp", "2", "-m", "1G",
		"-nographic",
		"-kernel", filepath.Join(dir, "Image"),
		"-drive", "file=" + filepath.Join(dir, "rootfs.ext4") + ",format=raw,if=virtio",
		"-snapshot",
		"-fsdev", "local,id=moddev,path=" + filepath.Join(dir, "modules") + ",security_model=none,readonly=on",
		"-device", "virtio-9p-pci,fsdev=moddev,mount_tag=modules_mount",
		"-netdev", "user,id=net0", "-device", "virtio-net-device,netdev=net0",
		"-append", "root=/dev/vda rw init=/init console=ttyAMA0",
	}
	if _, err := os.Stat(filepath.Join(dir, "virt-kdal.dtb")); err == nil {
		args = append(args, "-dtb", filepath.Join(dir, "virt-kdal.dtb"))
	}
	return args
}
