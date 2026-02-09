// Package builder provides kernel and module build orchestration for elmos.
package builder

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	pb "github.com/NguyenTrongPhuc552003/elmos/api/proto"
	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
)

// KernelBuilder orchestrates kernel build operations.
type KernelBuilder struct {
	exec executor.Executor
	fs   filesystem.FileSystem
	cfg  *elconfig.Config
	ctx  *elcontext.Context
	tm   *toolchain.Manager
}

// NewKernelBuilder creates a new KernelBuilder with the given dependencies.
func NewKernelBuilder(exec executor.Executor, fs filesystem.FileSystem, cfg *elconfig.Config, ctx *elcontext.Context, tm *toolchain.Manager) *KernelBuilder {
	return &KernelBuilder{
		exec: exec,
		fs:   fs,
		cfg:  cfg,
		ctx:  ctx,
		tm:   tm,
	}
}

// BuildOptions contains options for building the kernel.
type BuildOptions struct {
	Jobs    int
	Targets []string
}

// Build builds the kernel with the specified targets.
func (b *KernelBuilder) Build(ctx context.Context, opts BuildOptions) error {
	// Validate targets
	for _, target := range opts.Targets {
		if !elconfig.ValidBuildTargets[target] {
			return fmt.Errorf("invalid build target: %s", target)
		}
	}

	// Determine job count
	jobs := opts.Jobs
	if jobs <= 0 {
		jobs = b.cfg.Build.Jobs
	}

	// Get environment with correct toolchain
	env, crossCompile, err := getToolchainEnv(b.ctx, b.cfg, b.tm, b.fs, b.cfg.Build.Arch)
	if err != nil {
		return fmt.Errorf("failed to configure toolchain environment: %w", err)
	}

	// Build make arguments
	args := []string{
		"-C", b.cfg.Paths.KernelDir,
		fmt.Sprintf("-j%d", jobs),
		fmt.Sprintf("ARCH=%s", b.cfg.Build.Arch),
		"LLVM=1",
		fmt.Sprintf("CROSS_COMPILE=%s", crossCompile),
	}
	args = append(args, opts.Targets...)

	return b.exec.RunWithEnv(ctx, env, "make", args...)
}

// BuildWithProgress builds the kernel with streaming progress updates.
// Returns a channel that sends BuildProgress events (stage, log, error, complete).
func (b *KernelBuilder) BuildWithProgress(ctx context.Context, opts BuildOptions) (<-chan *pb.BuildProgress, error) {
	// Validate targets
	for _, target := range opts.Targets {
		if !elconfig.ValidBuildTargets[target] {
			return nil, fmt.Errorf("invalid build target: %s", target)
		}
	}

	// Determine job count
	jobs := opts.Jobs
	if jobs <= 0 {
		jobs = b.cfg.Build.Jobs
	}

	// Get environment with correct toolchain
	env, crossCompile, err := getToolchainEnv(b.ctx, b.cfg, b.tm, b.fs, b.cfg.Build.Arch)
	if err != nil {
		return nil, fmt.Errorf("failed to configure toolchain environment: %w", err)
	}

	// Build make arguments
	args := []string{
		"-C", b.cfg.Paths.KernelDir,
		fmt.Sprintf("-j%d", jobs),
		fmt.Sprintf("ARCH=%s", b.cfg.Build.Arch),
		"LLVM=1",
		fmt.Sprintf("CROSS_COMPILE=%s", crossCompile),
	}
	args = append(args, opts.Targets...)

	// Create progress channel
	progressCh := make(chan *pb.BuildProgress, 100)

	go func() {
		defer close(progressCh)

		// Send initial stage
		progressCh <- &pb.BuildProgress{
			Event: &pb.BuildProgress_Stage{
				Stage: &pb.BuildStage{
					Name:      "Starting",
					Progress:  0,
					Component: fmt.Sprintf("Building targets: %s", strings.Join(opts.Targets, ", ")),
				},
			},
		}

		// Start streaming command
		linesCh, errCh := b.exec.RunWithEnvStreaming(ctx, env, "make", args...)

		// Pattern to detect build progress (e.g., "CC kernel/fork.o", "LD vmlinux")
		buildPattern := regexp.MustCompile(`^\s*(CC|LD|AR|AS|OBJCOPY|GEN|CHK|UPD)\s+(.+)$`)

		var linesProcessed int
		var lastProgress int32

		// Process streaming output
		for {
			select {
			case line, ok := <-linesCh:
				if !ok {
					// No more lines, wait for error channel
					goto waitForError
				}

				linesProcessed++

				// Send log line
				progressCh <- &pb.BuildProgress{
					Event: &pb.BuildProgress_Log{
						Log: &pb.LogLine{
							Message: line,
							Level:   pb.LogLine_INFO,
						},
					},
				}

				// Parse build progress
				if matches := buildPattern.FindStringSubmatch(line); matches != nil {
					stage := matches[1]
					file := matches[2]

					// Estimate progress based on stage and line count
					// This is a rough heuristic - real implementation would track file count
					var progress int32
					switch stage {
					case "CC", "AS":
						progress = int32(float32(linesProcessed) / 100.0 * 70) // Compilation is 70%
					case "LD":
						progress = 80
					case "AR":
						progress = 85
					case "OBJCOPY":
						progress = 95
					}

					if progress > lastProgress && progress <= 100 {
						lastProgress = progress
						progressCh <- &pb.BuildProgress{
							Event: &pb.BuildProgress_Stage{
								Stage: &pb.BuildStage{
									Name:      stage,
									Progress:  progress,
									Component: fmt.Sprintf("Processing %s", file),
								},
							},
						}
					}
				}

			case <-ctx.Done():
				progressCh <- &pb.BuildProgress{
					Event: &pb.BuildProgress_Error{
						Error: &pb.BuildError{
							Message: "Build cancelled",
						},
					},
				}
				return
			}
		}

	waitForError:
		// Wait for command completion
		if err := <-errCh; err != nil {
			progressCh <- &pb.BuildProgress{
				Event: &pb.BuildProgress_Error{
					Error: &pb.BuildError{
						Message: err.Error(),
					},
				},
			}
		} else {
			// Success
			progressCh <- &pb.BuildProgress{
				Event: &pb.BuildProgress_Complete{
					Complete: &pb.BuildComplete{
						Success:    true,
						DurationMs: 0, // TODO: track time
						ImagePath:  b.cfg.Paths.KernelDir,
					},
				},
			}
		}
	}()

	return progressCh, nil
}

// Configure runs kernel configuration (menuconfig, defconfig, etc.).
func (b *KernelBuilder) Configure(ctx context.Context, configType string) error {
	// Validate config type
	valid := false
	for _, ct := range elconfig.KernelConfigTypes {
		if ct == configType {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid config type: %s", configType)
	}

	// Get environment with correct toolchain
	env, crossCompile, err := getToolchainEnv(b.ctx, b.cfg, b.tm, b.fs, b.cfg.Build.Arch)
	if err != nil {
		return fmt.Errorf("failed to configure toolchain environment: %w", err)
	}

	args := []string{
		"-C", b.cfg.Paths.KernelDir,
		fmt.Sprintf("ARCH=%s", b.cfg.Build.Arch),
		"LLVM=1",
		fmt.Sprintf("CROSS_COMPILE=%s", crossCompile),
		configType,
	}

	if err := b.exec.RunWithEnv(ctx, env, "make", args...); err != nil {
		return err
	}

	// For kvm_guest.config, force graphics options to be built-in (=y)
	if configType == "kvm_guest.config" {
		if err := b.forceGraphicsConfig(ctx); err != nil {
			return fmt.Errorf("failed to enable graphics options: %w", err)
		}
	}

	return nil
}

// forceGraphicsConfig forces graphics-related options to be built-in for QEMU GUI.
func (b *KernelBuilder) forceGraphicsConfig(ctx context.Context) error {
	scriptPath := filepath.Join(b.cfg.Paths.KernelDir, "scripts", "config")
	configFile := filepath.Join(b.cfg.Paths.KernelDir, ".config")

	// Use scripts/config to force graphics options to =y (built-in)
	options := []string{
		"--enable", "CONFIG_DRM",
		"--enable", "CONFIG_DRM_VIRTIO_GPU",
		"--enable", "CONFIG_FB",
		"--enable", "CONFIG_FRAMEBUFFER_CONSOLE",
	}

	configArgs := append([]string{"--file", configFile}, options...)
	if err := b.exec.Run(ctx, scriptPath, configArgs...); err != nil {
		return fmt.Errorf("scripts/config failed: %w", err)
	}

	// Run olddefconfig to finalize the config
	args := []string{
		"-C", b.cfg.Paths.KernelDir,
		fmt.Sprintf("ARCH=%s", b.cfg.Build.Arch),
		"LLVM=1",
		"olddefconfig",
	}

	env := b.ctx.GetMakeEnv()
	return b.exec.RunWithEnv(ctx, env, "make", args...)
}

// EnableKVMConfig enables KVM-specific kernel config options.
func (b *KernelBuilder) EnableKVMConfig(ctx context.Context) error {
	configFile := filepath.Join(b.cfg.Paths.KernelDir, ".config")

	if !b.fs.Exists(configFile) {
		return fmt.Errorf(".config not found - run 'elmos kernel config' first")
	}

	// Read current config
	content, err := b.fs.ReadFile(configFile)
	if err != nil {
		return err
	}

	// This is a simplified implementation
	// In practice, you'd use scripts/config or similar
	_ = content // For now, just run olddefconfig which applies defaults

	args := []string{
		"-C", b.cfg.Paths.KernelDir,
		fmt.Sprintf("ARCH=%s", b.cfg.Build.Arch),
		"LLVM=1",
		"olddefconfig",
	}

	env := b.ctx.GetMakeEnv()
	return b.exec.RunWithEnv(ctx, env, "make", args...)
}

// Clean runs distclean on the kernel source.
func (b *KernelBuilder) Clean(ctx context.Context) error {
	args := []string{
		"-C", b.cfg.Paths.KernelDir,
		fmt.Sprintf("ARCH=%s", b.cfg.Build.Arch),
		"LLVM=1",
		"distclean",
	}

	env := b.ctx.GetMakeEnv()
	return b.exec.RunWithEnvSilent(ctx, env, "make", args...)
}

// GetDefaultTargets returns the default build targets for the current architecture.
func (b *KernelBuilder) GetDefaultTargets() []string {
	return b.ctx.GetDefaultTargets()
}

// HasConfig checks if the kernel has been configured.
func (b *KernelBuilder) HasConfig() bool {
	return b.ctx.HasConfig()
}

// HasKernelImage checks if the kernel image has been built.
func (b *KernelBuilder) HasKernelImage() bool {
	return b.ctx.HasKernelImage()
}
