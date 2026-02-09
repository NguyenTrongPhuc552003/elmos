// Package server implements gRPC service handlers for ELMOS.
package server

import (
	"fmt"
	"net"

	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	elcontext "github.com/NguyenTrongPhuc552003/elmos/core/context"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/builder"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/emulator"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/toolchain"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/executor"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/filesystem"
	"github.com/NguyenTrongPhuc552003/elmos/core/infra/volume"
	"google.golang.org/grpc"

	pb "github.com/NguyenTrongPhuc552003/elmos/api/proto"
)

// Server is the main gRPC server for ELMOS.
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener

	// Services
	workspace *WorkspaceServer
	kernel    *KernelServer
	toolchain *ToolchainServer
	qemu      *QEMUServer
}

// Config contains configuration for the gRPC server.
type Config struct {
	// Address is the listen address (e.g., "unix:///tmp/elmos.sock" or "tcp://localhost:50051")
	Address string
}

// NewServer creates a new gRPC server with all services registered.
func NewServer(cfg *elconfig.Config, ctx *elcontext.Context, exec executor.Executor, fs filesystem.FileSystem, volMgr volume.Manager) (*Server, error) {
	// Create domain services
	tm := toolchain.NewManager(exec, fs, cfg, nil) // Pass nil for printer (server doesn't need it)
	kb := builder.NewKernelBuilder(exec, fs, cfg, ctx, tm)
	runner := emulator.NewQEMURunner(exec, fs, cfg, ctx)

	// Create gRPC servers for each service
	workspaceServer := NewWorkspaceServer(cfg, volMgr)
	kernelServer := NewKernelServer(cfg, kb)
	toolchainServer := NewToolchainServer(cfg, tm)
	qemuServer := NewQEMUServer(cfg, runner)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterWorkspaceServiceServer(grpcServer, workspaceServer)
	pb.RegisterKernelServiceServer(grpcServer, kernelServer)
	pb.RegisterToolchainServiceServer(grpcServer, toolchainServer)
	pb.RegisterQEMUServiceServer(grpcServer, qemuServer)

	return &Server{
		grpcServer: grpcServer,
		workspace:  workspaceServer,
		kernel:     kernelServer,
		toolchain:  toolchainServer,
		qemu:       qemuServer,
	}, nil
}

// Serve starts the gRPC server on the specified address.
func (s *Server) Serve(address string) error {
	// Parse address (unix:// or tcp://)
	network, addr, err := parseAddress(address)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	// Create listener
	listener, err := net.Listen(network, addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	s.listener = listener

	// Start serving
	return s.grpcServer.Serve(listener)
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// parseAddress parses a gRPC address in the format "unix:///path" or "tcp://host:port".
func parseAddress(address string) (network string, addr string, err error) {
	if len(address) > 7 && address[:7] == "unix://" {
		return "unix", address[7:], nil
	}
	if len(address) > 6 && address[:6] == "tcp://" {
		return "tcp", address[6:], nil
	}
	// Default to TCP if no prefix
	return "tcp", address, nil
}
