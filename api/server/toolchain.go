// Package server implements gRPC service handlers for ELMOS.
package server

import (
	"context"
	"fmt"

	pb "github.com/NguyenTrongPhuc552003/elmos/api/proto"
	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
)

// ToolchainServer implements the ToolchainService gRPC service.
type ToolchainServer struct {
	pb.UnimplementedToolchainServiceServer
	cfg *elconfig.Config
	tm  *toolchain.Manager
}

// NewToolchainServer creates a new ToolchainServer instance.
func NewToolchainServer(cfg *elconfig.Config, tm *toolchain.Manager) *ToolchainServer {
	return &ToolchainServer{
		cfg: cfg,
		tm:  tm,
	}
}

// Install installs crosstool-ng with streaming progress.
func (s *ToolchainServer) Install(req *pb.InstallRequest, stream pb.ToolchainService_InstallServer) error {
	ctx := stream.Context()

	// Send initial progress
	if err := stream.Send(&pb.InstallProgress{
		Stage:   pb.InstallProgress_DOWNLOADING,
		Message: "Installing crosstool-ng",
	}); err != nil {
		return err
	}

	// Install crosstool-ng
	if err := s.tm.Install(ctx); err != nil {
		return fmt.Errorf("failed to install crosstool-ng: %w", err)
	}

	// Send completion
	if err := stream.Send(&pb.InstallProgress{
		Stage:   pb.InstallProgress_COMPLETE,
		Message: "crosstool-ng installed successfully",
	}); err != nil {
		return err
	}

	return nil
}

// ListSamples lists available toolchain samples.
func (s *ToolchainServer) ListSamples(req *pb.ListSamplesRequest, stream pb.ToolchainService_ListSamplesServer) error {
	ctx := context.Background()

	// Get available samples
	samples, err := s.tm.ListSamples(ctx)
	if err != nil {
		return fmt.Errorf("failed to list samples: %w", err)
	}

	// Stream each sample
	for _, sample := range samples {
		if err := stream.Send(&pb.SampleTarget{
			Name:        sample,
			Description: fmt.Sprintf("Toolchain target: %s", sample),
		}); err != nil {
			return err
		}
	}

	return nil
}

// SelectTarget configures the toolchain for a specific target.
func (s *ToolchainServer) SelectTarget(ctx context.Context, req *pb.SelectTargetRequest) (*pb.SelectTargetResponse, error) {
	// Configure toolchain target
	if err := s.tm.SelectTarget(ctx, req.Target); err != nil {
		return &pb.SelectTargetResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("failed to configure target: %v", err),
		}, nil
	}

	return &pb.SelectTargetResponse{
		Success: true,
	}, nil
}

// Build builds the toolchain with streaming progress.
func (s *ToolchainServer) Build(req *pb.BuildToolchainRequest, stream pb.ToolchainService_BuildServer) error {
	ctx := stream.Context()

	// Send initial progress
	if err := stream.Send(&pb.BuildProgress{
		Event: &pb.BuildProgress_Stage{
			Stage: &pb.BuildStage{
				Name:      "Starting",
				Progress:  0,
				Component: "Toolchain build",
			},
		},
	}); err != nil {
		return err
	}

	// Build toolchain (Manager.Build requires jobs parameter)
	jobs := int(req.Jobs)
	if jobs <= 0 {
		jobs = 8 // default
	}

	if err := s.tm.Build(ctx, jobs); err != nil {
		return fmt.Errorf("toolchain build failed: %w", err)
	}

	// Send completion
	if err := stream.Send(&pb.BuildProgress{
		Event: &pb.BuildProgress_Complete{
			Complete: &pb.BuildComplete{
				Success:    true,
				DurationMs: 0, // TODO: track time
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

// GetStatus returns the toolchain build status.
func (s *ToolchainServer) GetStatus(ctx context.Context, req *pb.ToolchainStatusRequest) (*pb.ToolchainStatusResponse, error) {
	isInstalled := s.tm.IsInstalled()

	// Get installed toolchains
	toolchains, err := s.tm.GetInstalledToolchains()
	if err != nil {
		toolchains = nil
	}

	// Convert to proto format
	var protoToolchains []*pb.ToolchainInfo
	for _, tc := range toolchains {
		protoToolchains = append(protoToolchains, &pb.ToolchainInfo{
			Target:     tc.Target,
			GccVersion: tc.Version,
			Path:       tc.Path,
		})
	}

	return &pb.ToolchainStatusResponse{
		CrosstoolNgInstalled: isInstalled,
		CrosstoolNgVersion:   "unknown",
		Installed:            protoToolchains,
		SelectedTarget:       "",
	}, nil
}

// Clean cleans the toolchain build artifacts.
func (s *ToolchainServer) Clean(ctx context.Context, req *pb.CleanRequest) (*pb.CleanResponse, error) {
	// TODO: Implement toolchain clean
	return &pb.CleanResponse{
		Success: true,
		Message: "Toolchain artifacts cleaned",
	}, nil
}
