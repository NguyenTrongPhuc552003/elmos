package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/NguyenTrongPhuc552003/elmos/core/domain/demo"
)

// BuildDemo creates the demo command: prebuilt quick-start demos that need
// only QEMU (audit E6, quick-start levels 1 and 2).
func BuildDemo(ctx *Context) *cobra.Command {
	var cacheDir string
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Run a prebuilt demo in QEMU (no kernel build, no sudo)",
		Long: `Run a prebuilt demo. Assets are downloaded from the project's GitHub
releases, checked against the release's SHA256SUMS and cached, so later
runs work offline. Needs only qemu-system-aarch64 (and python3 and dtc
for the HAVEN demo).

  elmos demo linux   Linux v6.18 arm64 with an alpine:3 root file system;
                     the KDAL module loads at boot. Exit QEMU with Ctrl-A X.
  elmos demo haven   HAVEN: two partitions, then read-back finds a fault
                     that the partitions' own checks miss (about 15 s).`,
	}
	cmd.PersistentFlags().StringVar(&cacheDir, "cache", "", "Cache directory (default ~/.elmos/demo)")

	fetcher := func() (*demo.Fetcher, error) {
		dir := cacheDir
		if dir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			dir = filepath.Join(home, ".elmos", "demo")
		}
		return demo.NewFetcher(dir), nil
	}

	linuxCmd := &cobra.Command{
		Use:   "linux",
		Short: "Boot Linux v6.18 (arm64) with the KDAL module",
		RunE: func(c *cobra.Command, args []string) error {
			f, err := fetcher()
			if err != nil {
				return err
			}
			ctx.Printer.Step("Fetching %s (cached after the first run)...", demo.LinuxTag)
			dir, err := f.Release(c.Context(), "elmos", demo.LinuxTag)
			if err != nil {
				return err
			}
			if err := demo.Gunzip(filepath.Join(dir, "Image.gz"), filepath.Join(dir, "Image")); err != nil {
				return err
			}
			if err := demo.Gunzip(filepath.Join(dir, "rootfs.ext4.gz"), filepath.Join(dir, "rootfs.ext4")); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Join(dir, "modules"), 0o755); err != nil {
				return err
			}
			if err := demo.Untar(filepath.Join(dir, "modules.tar.gz"), filepath.Join(dir, "modules")); err != nil {
				return err
			}
			ctx.Printer.Step("Booting (exit with Ctrl-A X)...")
			return ctx.Exec.Run(c.Context(), "qemu-system-aarch64", demo.LinuxArgs(dir)...)
		},
	}

	havenCmd := &cobra.Command{
		Use:   "haven",
		Short: "Run the HAVEN try bundle (two partitions, read-back)",
		RunE: func(c *cobra.Command, args []string) error {
			f, err := fetcher()
			if err != nil {
				return err
			}
			ctx.Printer.Step("Fetching HAVEN %s (cached after the first run)...", demo.HavenTag)
			dir, err := f.Release(c.Context(), "haven", demo.HavenTag)
			if err != nil {
				return err
			}
			if err := demo.Untar(filepath.Join(dir, demo.HavenBundle), dir); err != nil {
				return err
			}
			bundle := filepath.Join(dir, trimTarGz(demo.HavenBundle))
			if _, err := os.Stat(filepath.Join(bundle, "try.sh")); err != nil {
				return fmt.Errorf("bundle has no try.sh: %w", err)
			}
			return ctx.Exec.RunInDir(c.Context(), bundle, "sh", "./try.sh")
		},
	}

	cmd.AddCommand(linuxCmd, havenCmd)
	return cmd
}

func trimTarGz(name string) string {
	const ext = ".tar.gz"
	if len(name) > len(ext) && name[len(name)-len(ext):] == ext {
		return name[:len(name)-len(ext)]
	}
	return name
}
