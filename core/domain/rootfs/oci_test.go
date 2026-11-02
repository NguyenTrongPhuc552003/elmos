package rootfs

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"codeberg.org/NguyenTrongPhuc552003/elmos/assets"
	elconfig "codeberg.org/NguyenTrongPhuc552003/elmos/core/config"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"codeberg.org/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
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

// The default engine is OCI: it writes the image without sudo, then turns
// the compact read-only ext4 into a writable one of the asked size.
func TestCreateDefaultsToOCIAndFixesUpTheImage(t *testing.T) {
	dir := t.TempDir()
	cfg := &elconfig.Config{}
	cfg.Build.Arch = "arm64"
	cfg.Paths.DiskImage = filepath.Join(dir, "disk.img")
	cfg.Paths.RootfsDir = filepath.Join(dir, "rootfs")
	rec := executor.NewRecorder()
	rec.Outputs["brew --prefix e2fsprogs"] = "/brew/e2fsprogs\n"
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
	fi, err := os.Stat(cfg.Paths.DiskImage)
	if err != nil || fi.Size() == 0 {
		t.Fatalf("disk image not written: %v", err)
	}
	var tools []string
	for _, call := range rec.Calls {
		if strings.HasPrefix(call, "sudo") {
			t.Fatalf("OCI engine ran sudo: %s", call)
		}
		if strings.HasPrefix(call, "/brew/e2fsprogs/sbin/") {
			tools = append(tools, strings.Fields(strings.TrimPrefix(call, "/brew/e2fsprogs/sbin/"))[0])
		}
	}
	want := "e2fsck tune2fs resize2fs tune2fs e2fsck"
	if strings.Join(tools, " ") != want {
		t.Fatalf("fix-up = %v, want %s", tools, want)
	}
	if _, err := assets.GetInitScript(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateRejectsUnknownEngine(t *testing.T) {
	cfg := &elconfig.Config{}
	c := NewCreator(executor.NewRecorder(), filesystem.NewOSFileSystem(), cfg)
	if err := c.Create(context.Background(), CreateOptions{Engine: "nix"}); err == nil {
		t.Fatal("unknown engine accepted")
	}
}
