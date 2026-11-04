package demo

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sum(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// release serves files under /<repo>/releases/download/<tag>/ and counts requests.
func release(t *testing.T, files map[string][]byte, sums string) (*httptest.Server, *int) {
	t.Helper()
	n := new(int)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*n++
		name := filepath.Base(r.URL.Path)
		if name == "SHA256SUMS" {
			fmt.Fprint(w, sums)
			return
		}
		b, ok := files[name]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Write(b)
	}))
	t.Cleanup(srv.Close)
	return srv, n
}

func TestReleaseVerifiesAndCaches(t *testing.T) {
	files := map[string][]byte{"Image.gz": []byte("kernel"), "rootfs.ext4.gz": []byte("rootfs")}
	sums := sum(files["Image.gz"]) + "  Image.gz\n" + sum(files["rootfs.ext4.gz"]) + "  rootfs.ext4.gz\n"
	srv, n := release(t, files, sums)
	f := &Fetcher{Client: srv.Client(), CacheDir: t.TempDir(), Base: srv.URL}

	dir, err := f.Release(context.Background(), "elmos", "demo-x")
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "Image.gz")); string(b) != "kernel" {
		t.Fatalf("Image.gz = %q", b)
	}
	first := *n
	if _, err := f.Release(context.Background(), "elmos", "demo-x"); err != nil {
		t.Fatal(err)
	}
	if *n != first+1 { // only SHA256SUMS again
		t.Fatalf("second run made %d requests, want 1", *n-first)
	}
}

func TestReleaseFallsBackToMirror(t *testing.T) {
	files := map[string][]byte{"Image.gz": []byte("kernel")}
	srv, _ := release(t, files, sum(files["Image.gz"])+"  Image.gz\n")
	missing := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(missing.Close)
	f := &Fetcher{Client: srv.Client(), CacheDir: t.TempDir(), Base: missing.URL, Mirrors: []string{srv.URL}}
	dir, err := f.Release(context.Background(), "elmos", "demo-x")
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "Image.gz")); string(b) != "kernel" {
		t.Fatalf("Image.gz = %q", b)
	}
}

func TestReleaseRejectsWrongHash(t *testing.T) {
	files := map[string][]byte{"Image.gz": []byte("tampered")}
	srv, _ := release(t, files, sum([]byte("kernel"))+"  Image.gz\n")
	f := &Fetcher{Client: srv.Client(), CacheDir: t.TempDir(), Base: srv.URL}
	_, err := f.Release(context.Background(), "elmos", "demo-x")
	if err == nil || !strings.Contains(err.Error(), "sha256") {
		t.Fatalf("err = %v, want a hash mismatch", err)
	}
	if _, statErr := os.Stat(filepath.Join(f.CacheDir, "elmos", "demo-x", "Image.gz")); !os.IsNotExist(statErr) {
		t.Fatal("tampered file kept in the cache")
	}
}

func TestReleaseRejectsUnsafeNames(t *testing.T) {
	srv, _ := release(t, nil, sum([]byte("x"))+"  ../escape\n")
	f := &Fetcher{Client: srv.Client(), CacheDir: t.TempDir(), Base: srv.URL}
	if _, err := f.Release(context.Background(), "elmos", "demo-x"); err == nil {
		t.Fatal("path traversal in SHA256SUMS accepted")
	}
}

func TestUntarRefusesEscape(t *testing.T) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(zw)
	tw.WriteHeader(&tar.Header{Name: "../evil", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg})
	tw.Write([]byte("x"))
	tw.Close()
	zw.Close()
	dir := t.TempDir()
	src := filepath.Join(dir, "a.tar.gz")
	os.WriteFile(src, buf.Bytes(), 0o644)
	if err := Untar(src, filepath.Join(dir, "out")); err == nil {
		t.Fatal("archive escaping its directory accepted")
	}
}

func TestLinuxArgsUseShippedDTB(t *testing.T) {
	dir := t.TempDir()
	if strings.Contains(strings.Join(LinuxArgs(dir), " "), "-dtb") {
		t.Fatal("-dtb without a shipped device tree")
	}
	os.WriteFile(filepath.Join(dir, "virt-kdal.dtb"), []byte("x"), 0o644)
	if !strings.Contains(strings.Join(LinuxArgs(dir), " "), "-dtb "+filepath.Join(dir, "virt-kdal.dtb")) {
		t.Fatal("shipped device tree not used")
	}
}

func TestLinuxArgsKeepCacheClean(t *testing.T) {
	c := filepath.Join(string(os.PathSeparator), "c")
	args := strings.Join(LinuxArgs(c), " ")
	for _, want := range []string{"-kernel " + filepath.Join(c, "Image"),
		"file=" + filepath.Join(c, "rootfs.ext4") + ",format=raw", "-snapshot",
		"path=" + filepath.Join(c, "modules"), "readonly=on", "mount_tag=modules_mount"} {
		if !strings.Contains(args, want) {
			t.Fatalf("args lack %q: %s", want, args)
		}
	}
}
