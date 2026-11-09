// Command rootfs writes the SDK's root filesystem image with ELMOS's own
// OCI engine (the code behind 'elmos rootfs create'), to an explicit path.
//
// It exists for CI on Linux runners, where 'elmos init' (a macOS volume)
// cannot run: the creator only needs the two paths that a workspace would
// otherwise provide. See docs/developer/sdk-ci.md.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/rootfs"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func main() {
	out := flag.String("out", "rootfs.ext4", "disk image to write")
	image := flag.String("image", "alpine:3", "OCI image (linux/arm64)")
	size := flag.String("size", "1G", "disk image size")
	flag.Parse()

	abs, err := filepath.Abs(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	work, err := os.MkdirTemp("", "sdk-rootfs-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(work)

	cfg := &elconfig.Config{}
	cfg.Paths.DiskImage = abs
	cfg.Paths.RootfsDir = filepath.Join(work, "rootfs")
	c := rootfs.NewCreator(executor.NewShellExecutor(), filesystem.NewOSFileSystem(), cfg)
	opts := rootfs.CreateOptions{Engine: rootfs.EngineOCI, Image: *image, Size: *size}
	if err := c.Create(context.Background(), opts); err != nil {
		fmt.Fprintln(os.Stderr, "rootfs:", err)
		os.Exit(1)
	}
	fmt.Println("wrote", abs)
}
