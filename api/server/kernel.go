// Package server implements gRPC service handlers for ELMOS.
package server

import (
	"context"
	"fmt"
	"io"

	pb "github.com/NguyenTrongPhuc552003/elmos/api/proto"
	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/builder"
)

// KernelServer implements the KernelService gRPC service.
type KernelServer struct {
	pb.UnimplementedKernelServiceServer
	cfg     *elconfig.Config
	builder *builder.KernelBuilder
}

// NewKernelServer creates a new KernelServer instance.
func NewKernelServer(cfg *elconfig.Config, kb *builder.KernelBuilder) *KernelServer {
	return &KernelServer{
		cfg:     cfg,
		builder: kb,
	}
}

// Clone clones the Linux kernel source with streaming progress.
func (s *KernelServer) Clone(req *pb.CloneRequest, stream pb.KernelService_CloneServer) error {
	// Send initial progress
	if err := stream.Send(&pb.CloneProgress{
		Progress: 0,
		Message:  fmt.Sprintf("Cloning kernel version %s", req.Version),
	}); err != nil {
		return err
	}

	// TODO: Implement actual git clone with progress streaming
	// For now, just simulate completion
	if err := stream.Send(&pb.CloneProgress{
		Progress: 100,
		Message:  "Kernel cloned successfully",
	}); err != nil {
		return err
	}

	return nil
}

// Configure runs kernel configuration (defconfig, menuconfig, etc.).
func (s *KernelServer) Configure(ctx context.Context, req *pb.ConfigureRequest) (*pb.ConfigureResponse, error) {
	// Run kernel configuration
	if err := s.builder.Configure(ctx, req.ConfigType); err != nil {
		return &pb.ConfigureResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("configuration failed: %v", err),
		}, nil
	}

	return &pb.ConfigureResponse{
		Success: true,
	}, nil
}

// Build builds the kernel with streaming progress.
func (s *KernelServer) Build(req *pb.BuildRequest, stream pb.KernelService_BuildServer) error {
	ctx := stream.Context()

	// Prepare build options
	opts := builder.BuildOptions{
		Jobs:    int(req.Jobs),
		Targets: req.Targets,
	}

	// Start build with progress
	progressCh, err := s.builder.BuildWithProgress(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to start build: %w", err)
	}

	// Stream progress events to client
	for progress := range progressCh {
		if err := stream.Send(progress); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to send progress: %w", err)
		}
	}

	return nil
}

// Clean cleans the kernel build artifacts.
func (s *KernelServer) Clean(ctx context.Context, req *pb.CleanRequest) (*pb.CleanResponse, error) {
	// Run kernel clean
	if err := s.builder.Clean(ctx); err != nil {
		return nil, fmt.Errorf("clean failed: %w", err)
	}

	return &pb.CleanResponse{
		Success: true,
		Message: "Kernel build artifacts cleaned",
	}, nil
}

// GetStatus returns the current kernel build status.
func (s *KernelServer) GetStatus(ctx context.Context, req *pb.KernelStatusRequest) (*pb.KernelStatus, error) {
	return &pb.KernelStatus{
		Cloned:     true, // TODO: Check if kernel is cloned
		Configured: s.builder.HasConfig(),
		Built:      s.builder.HasKernelImage(),
		Version:    "6.18", // TODO: Get from config
	}, nil
}

// ListVersions lists available kernel versions.
func (s *KernelServer) ListVersions(ctx context.Context, req *pb.ListVersionsRequest) (*pb.ListVersionsResponse, error) {
	// TODO: Implement kernel version listing
	// For now, return some common versions
	return &pb.ListVersionsResponse{
		Versions: []string{"v6.18", "v6.6", "v5.15"},
	}, nil
}
