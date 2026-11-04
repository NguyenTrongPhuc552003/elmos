package rootfs

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/NguyenTrongPhuc552003/elmos/assets"
	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// imageTar builds a tiny image file system: a directory, a setuid binary
// owned by root, a symlink, a hard link and an /init that must be replaced.
func imageTar(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	add := func(h *tar.Header, body string) {
		h.Size = int64(len(body))
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	add(&tar.Header{Name: "usr/", Typeflag: tar.TypeDir, Mode: 0o755}, "")
	add(&tar.Header{Name: "usr/bin/", Typeflag: tar.TypeDir, Mode: 0o755}, "")
	add(&tar.Header{Name: "usr/bin/su", Typeflag: tar.TypeReg, Mode: 0o4755, Uid: 0, Gid: 0}, "su")
	add(&tar.Header{Name: "usr/bin/sh", Typeflag: tar.TypeSymlink, Linkname: "dash", Mode: 0o777}, "")
	add(&tar.Header{Name: "usr/bin/su2", Typeflag: tar.TypeLink, Linkname: "usr/bin/su"}, "")
	add(&tar.Header{Name: "init", Typeflag: tar.TypeReg, Mode: 0o755}, "old init")
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestWithInitReplacesInitAndKeepsMetadata(t *testing.T) {
	initScript := []byte("#!/bin/sh\necho elmos\n")
	var out bytes.Buffer
	if err := withInit(&out, bytes.NewReader(imageTar(t)), initScript); err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(&out)
	inits := 0
	seen := map[string]*tar.Header{}
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		seen[h.Name] = h
		if h.Name == "init" {
			inits++
			body, _ := io.ReadAll(tr)
			if !bytes.Equal(body, initScript) || h.Mode != 0o755 {
				t.Fatalf("init = %q mode %o", body, h.Mode)
			}
		}
	}
	if inits != 1 {
		t.Fatalf("%d init entries, want 1", inits)
	}
	if h := seen["usr/bin/su"]; h == nil || h.Mode != 0o4755 {
		t.Fatalf("setuid mode lost: %+v", h)
	}
	if h := seen["usr/bin/su2"]; h == nil || h.Typeflag != tar.TypeLink {
		t.Fatalf("hard link lost: %+v", h)
	}
}

func TestOCIPlatform(t *testing.T) {
	for arch, want := range map[string]string{"arm64": "arm64", "arm": "arm", "riscv": "riscv64"} {
		p, err := ociPlatform(arch)
		if err != nil || p.Architecture != want || p.OS != "linux" {
			t.Fatalf("%s -> %+v %v", arch, p, err)
		}
	}
	if _, err := ociPlatform("mips"); err == nil {
		t.Fatal("unknown arch accepted")
	}
}

// The default engine is OCI: it writes a writable ext4 of the asked size in
// pure Go, runs nothing on the host (no sudo, no e2fsprogs), and the result
// passes e2fsck when e2fsck is available.
func TestCreateDefaultsToOCIAndWritesACleanImage(t *testing.T) {
	dir := t.TempDir()
	cfg := &elconfig.Config{}
	cfg.Build.Arch = "arm64"
	cfg.Paths.DiskImage = filepath.Join(dir, "disk.img")
	cfg.Paths.RootfsDir = filepath.Join(dir, "rootfs")
	rec := executor.NewRecorder()
	c := NewCreator(rec, filesystem.NewOSFileSystem(), cfg)
	var gotRef, gotArch string
	c.pull = func(_ context.Context, ref string, p *v1.Platform) (io.ReadCloser, error) {
		gotRef, gotArch = ref, p.Architecture
		return io.NopCloser(bytes.NewReader(imageTar(t))), nil
	}

	if err := c.Create(context.Background(), CreateOptions{Size: "64M"}); err != nil {
		t.Fatal(err)
	}
	if gotRef != DefaultOCIImage || gotArch != "arm64" {
		t.Fatalf("pulled %s for %s", gotRef, gotArch)
	}
	if len(rec.Calls) != 0 {
		t.Fatalf("OCI engine ran host commands: %v", rec.Calls)
	}
	fi, err := os.Stat(cfg.Paths.DiskImage)
	if err != nil || fi.Size() != 64<<20 {
		t.Fatalf("disk image size %v, %v; want 64 MiB", fi, err)
	}
	if _, err := assets.GetInitScript(); err != nil {
		t.Fatal(err)
	}
	e2fsck := findE2fsck()
	if e2fsck == "" {
		t.Skip("e2fsck not installed: structure not checked")
	}
	out, err := exec.Command(e2fsck, "-fn", cfg.Paths.DiskImage).CombinedOutput()
	if err != nil {
		t.Fatalf("e2fsck -fn: %v\n%s", err, out)
	}
	if strings.Contains(string(out), "read-only") {
		t.Fatalf("image still read-only:\n%s", out)
	}
}

func findE2fsck() string {
	for _, p := range []string{"/opt/homebrew/opt/e2fsprogs/sbin/e2fsck", "/usr/local/opt/e2fsprogs/sbin/e2fsck"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	if p, err := exec.LookPath("e2fsck"); err == nil {
		return p
	}
	return ""
}

func TestCreateRejectsUnknownEngine(t *testing.T) {
	cfg := &elconfig.Config{}
	c := NewCreator(executor.NewRecorder(), filesystem.NewOSFileSystem(), cfg)
	if err := c.Create(context.Background(), CreateOptions{Engine: "nix"}); err == nil {
		t.Fatal("unknown engine accepted")
	}
}
