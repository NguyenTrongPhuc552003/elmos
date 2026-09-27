package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

func TestUpdateInitConfigWorkspacePaths(t *testing.T) {
	tests := []struct {
		name      string
		paths     string
		kernelDir string
		rootfsDir string
		diskImage string
	}{
		{
			name:      "computed defaults follow new volume",
			kernelDir: "/Volumes/renamed/linux",
			rootfsDir: "/Volumes/renamed/rootfs",
			diskImage: "/Volumes/renamed/disk.img",
		},
		{
			name:      "explicit custom paths stay unchanged",
			paths:     "  kernel_dir: /custom/kernel\n  rootfs_dir: /custom/rootfs\n  disk_image: /custom/disk.img\n",
			kernelDir: "/custom/kernel",
			rootfsDir: "/custom/rootfs",
			diskImage: "/custom/disk.img",
		},
		{
			name:      "explicit old defaults stay unchanged",
			paths:     "  kernel_dir: /Volumes/elmos/linux\n  rootfs_dir: /Volumes/elmos/rootfs\n  disk_image: /Volumes/elmos/disk.img\n",
			kernelDir: "/Volumes/elmos/linux",
			rootfsDir: "/Volumes/elmos/rootfs",
			diskImage: "/Volumes/elmos/disk.img",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			configPath := filepath.Join(root, "elmos.yaml")
			content := "image:\n  volume_name: elmos\npaths:\n  project_root: " + root + "\n" + tt.paths
			if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				t.Fatal(err)
			}
			ctx := &Context{Config: cfg, FS: filesystem.NewOSFileSystem()}
			if err := updateInitConfig(ctx, "renamed", cfg.Image.Size); err != nil {
				t.Fatal(err)
			}

			assertWorkspacePaths(t, cfg, tt.kernelDir, tt.rootfsDir, tt.diskImage)
			reloaded, err := config.Load(configPath)
			if err != nil {
				t.Fatal(err)
			}
			assertWorkspacePaths(t, reloaded, tt.kernelDir, tt.rootfsDir, tt.diskImage)
		})
	}
}

func assertWorkspacePaths(t *testing.T, cfg *config.Config, kernelDir, rootfsDir, diskImage string) {
	t.Helper()
	if cfg.Paths.KernelDir != kernelDir || cfg.Paths.RootfsDir != rootfsDir || cfg.Paths.DiskImage != diskImage {
		t.Fatalf("paths = (%q, %q, %q), want (%q, %q, %q)",
			cfg.Paths.KernelDir, cfg.Paths.RootfsDir, cfg.Paths.DiskImage,
			kernelDir, rootfsDir, diskImage)
	}
}
