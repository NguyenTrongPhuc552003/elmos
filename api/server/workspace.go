// Package server implements gRPC service handlers for ELMOS.
package server

import (
	"context"
	"fmt"

	pb "github.com/NguyenTrongPhuc552003/elmos/api/proto"
	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/volume"
)

// WorkspaceServer implements the WorkspaceService gRPC service.
type WorkspaceServer struct {
	pb.UnimplementedWorkspaceServiceServer
	cfg    *elconfig.Config
	volMgr volume.Manager
}

// NewWorkspaceServer creates a new WorkspaceServer instance.
func NewWorkspaceServer(cfg *elconfig.Config, volMgr volume.Manager) *WorkspaceServer {
	return &WorkspaceServer{
		cfg:    cfg,
		volMgr: volMgr,
	}
}

// Init initializes a new ELMOS workspace with streaming progress.
func (s *WorkspaceServer) Init(req *pb.InitRequest, stream pb.WorkspaceService_InitServer) error {
	ctx := stream.Context()

	// Send initial progress
	if err := stream.Send(&pb.InitProgress{
		Stage:   pb.InitProgress_CREATING_VOLUME,
		Message: fmt.Sprintf("Creating volume %s with size %s", req.Name, req.Size),
	}); err != nil {
		return err
	}

	// Create volume
	volumePath := "" // Volume path determined by manager
	if err := s.volMgr.Create(ctx, req.Name, req.Size, volumePath); err != nil {
		return fmt.Errorf("failed to create volume: %w", err)
	}

	// Send mount progress
	if err := stream.Send(&pb.InitProgress{
		Stage:   pb.InitProgress_MOUNTING,
		Message: "Mounting volume",
	}); err != nil {
		return err
	}

	// Mount volume
	mountPoint := "" // Mount point determined by manager
	if err := s.volMgr.Mount(ctx, volumePath, mountPoint); err != nil {
		return fmt.Errorf("failed to mount volume: %w", err)
	}

	// Send config progress
	if err := stream.Send(&pb.InitProgress{
		Stage:   pb.InitProgress_SAVING_CONFIG,
		Message: "Saving workspace configuration",
	}); err != nil {
		return err
	}

	// Send completion
	if err := stream.Send(&pb.InitProgress{
		Stage:   pb.InitProgress_COMPLETE,
		Message: "Workspace initialized successfully",
	}); err != nil {
		return err
	}

	return nil
}

// Mount mounts an existing workspace volume.
func (s *WorkspaceServer) Mount(ctx context.Context, req *pb.MountRequest) (*pb.MountResponse, error) {
	mountPoint := "" // Mount point determined by manager
	if err := s.volMgr.Mount(ctx, req.VolumePath, mountPoint); err != nil {
		return &pb.MountResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("failed to mount volume: %v", err),
		}, nil
	}

	return &pb.MountResponse{
		Success:    true,
		MountPoint: mountPoint,
	}, nil
}

// Unmount unmounts a workspace volume.
func (s *WorkspaceServer) Unmount(ctx context.Context, req *pb.UnmountRequest) (*pb.UnmountResponse, error) {
	if err := s.volMgr.Unmount(ctx, req.MountPoint, false); err != nil {
		return &pb.UnmountResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("failed to unmount volume: %v", err),
		}, nil
	}

	return &pb.UnmountResponse{
		Success: true,
	}, nil
}

// GetStatus returns the status of a workspace.
func (s *WorkspaceServer) GetStatus(ctx context.Context, req *pb.WorkspaceStatusRequest) (*pb.WorkspaceStatus, error) {
	// TODO: Implement actual workspace status check
	return &pb.WorkspaceStatus{
		Mounted:    false,
		MountPoint: "",
		VolumePath: "",
		SizeBytes:  0,
		UsedBytes:  0,
		AvailBytes: 0,
	}, nil
}
