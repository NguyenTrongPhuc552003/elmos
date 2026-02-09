package emulator

import (
	"fmt"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
)

// buildArgs constructs the QEMU command line arguments.
func (q *QEMURunner) buildArgs(archCfg *elconfig.ArchConfig, kernelImage string, opts RunOptions) []string {
	args := []string{
		"-m", q.cfg.QEMU.Memory,
		"-smp", fmt.Sprintf("%d", q.cfg.QEMU.SMP),
		"-kernel", kernelImage,
		"-machine", archCfg.QEMUMachine,
	}

	if archCfg.QEMUCPU != "" {
		args = append(args, "-cpu", archCfg.QEMUCPU)
	}

	if archCfg.QEMUBios != "" {
		args = append(args, "-bios", "default")
	}

	// Disk and networking
	args = append(args,
		"-drive", fmt.Sprintf("file=%s,format=raw,if=virtio", q.cfg.Paths.DiskImage),
		"-device", "virtio-net-device,netdev=net0",
		"-netdev", fmt.Sprintf("user,id=net0,hostfwd=tcp::%d-:22", q.cfg.QEMU.SSHPort),
	)

	// 9p share for modules
	args = append(args,
		"-fsdev", fmt.Sprintf("local,id=moddev,path=%s,security_model=none", q.cfg.Paths.ModulesDir),
		"-device", "virtio-9p-pci,fsdev=moddev,mount_tag=modules_mount",
	)

	// Boot parameters
	appendStr := "root=/dev/vda rw init=/init earlycon"

	// Display mode
	if opts.Graphical {
		args = append(args, "-display", "cocoa")
		args = append(args,
			"-device", "virtio-gpu-pci",
			"-device", "virtio-keyboard-pci",
			"-device", "virtio-mouse-pci",
		)
		appendStr += " console=tty0"
	} else {
		args = append(args,
			"-nographic",
			"-serial", "mon:stdio",
		)
		appendStr += fmt.Sprintf(" console=%s", archCfg.Console)
	}

	args = append(args, "-append", appendStr)

	// Debug flags
	if opts.Debug {
		args = append(args, "-s", "-S")
	}

	return args
}
