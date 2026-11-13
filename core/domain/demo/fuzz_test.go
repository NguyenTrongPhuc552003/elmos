package demo

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tarOf builds an uncompressed tar with one regular file per name.
func tarOf(t testing.TB, names ...string) []byte {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, n := range names {
		if err := tw.WriteHeader(&tar.Header{Name: n, Mode: 0o644, Size: 1, Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte("x")); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// FuzzUntar feeds Untar any tar stream (gzip-wrapped here, so the fuzzer
// works on the tar structure) and checks that nothing appears outside the
// target directory: Untar unpacks releases downloaded from the network.
func FuzzUntar(f *testing.F) {
	for _, names := range [][]string{{"a.txt"}, {"d/b.txt"}, {"../escape"}, {"/abs"}, {"d/../../escape"}} {
		f.Add(tarOf(f, names...))
	}
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, data []byte) {
		root := t.TempDir()
		dir := filepath.Join(root, "out")
		src := filepath.Join(root, "in.tar.gz")
		var gz bytes.Buffer
		zw := gzip.NewWriter(&gz)
		if _, err := zw.Write(data); err != nil {
			t.Fatal(err)
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(src, gz.Bytes(), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		_ = Untar(src, dir) // errors are fine; escapes are not
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			if e.Name() != "out" && e.Name() != "in.tar.gz" {
				t.Fatalf("Untar wrote outside its directory: %s", e.Name())
			}
		}
		err = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Type()&os.ModeSymlink != 0 {
				t.Fatalf("Untar created a symbolic link: %s", p)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}

// FuzzParseSums checks that every name parseSums accepts is a plain file
// name (Release joins it to its cache directory) and that both results
// agree.
func FuzzParseSums(f *testing.F) {
	h := strings.Repeat("a", 64)
	for _, s := range []string{h + "  Image.gz\n", h + " *rootfs.ext4\n", h + "  ../x\n", "garbage\n", ""} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, text string) {
		p := filepath.Join(t.TempDir(), "SHA256SUMS")
		if err := os.WriteFile(p, []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
		sums, order, err := parseSums(p)
		if err != nil {
			return
		}
		for _, n := range order {
			if n != filepath.Base(n) || n == "." || n == ".." || strings.ContainsAny(n, `/\`) {
				t.Fatalf("unsafe name accepted: %q", n)
			}
			if _, ok := sums[n]; !ok {
				t.Fatalf("%q in order but not in sums", n)
			}
		}
	})
}
