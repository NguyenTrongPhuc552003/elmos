// Package server implements gRPC service handlers for ELMOS.
package server

import (
	"context"

	pb "github.com/NguyenTrongPhuc552003/elmos/api/proto"
	elconfig "github.com/NguyenTrongPhuc552003/elmos/core/config"
	"github.com/NguyenTrongPhuc552003/elmos/core/domain/emulator"
)

// QEMUServer implements the QEMUService gRPC service.
type QEMUServer struct {
	pb.UnimplementedQEMUServiceServer
	cfg    *elconfig.Config
	runner *emulator.QEMURunner
}

// NewQEMUServer creates a new QEMUServer instance.
func NewQEMUServer(cfg *elconfig.Config, runner *emulator.QEMURunner) *QEMUServer {
	return &QEMUServer{
		cfg:    cfg,
		runner: runner,
	}
}

// Run starts QEMU and streams console output.
func (s *QEMUServer) Run(req *pb.QEMURunRequest, stream pb.QEMUService_RunServer) error {
	ctx := stream.Context()

	// Send initial started event
	if err := stream.Send(&pb.QEMUOutput{
		Event: &pb.QEMUOutput_Started{
			Started: &pb.QEMUStarted{
				Pid:         0, // TODO: Get actual PID
				QemuVersion: "QEMU emulator",
				Command:     "qemu-system-...",
			},
		},
	}); err != nil {
		return err
	}

	// TODO: Implement actual QEMU execution with console streaming
	// For now, simulate console output
	if err := stream.Send(&pb.QEMUOutput{
		Event: &pb.QEMUOutput_Console{
			Console: &pb.ConsoleOutput{
				Data:        []byte("Booting Linux...\n"),
				TimestampMs: 0,
			},
		},
	}); err != nil {
		return err
	}

	// Keep stream alive until context is cancelled
	<-ctx.Done()

	// Send stopped event
	if err := stream.Send(&pb.QEMUOutput{
		Event: &pb.QEMUOutput_Stopped{
			Stopped: &pb.QEMUStopped{
				ExitCode: 0,
				UptimeMs: 0,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

// Stop stops a running QEMU instance.
func (s *QEMUServer) Stop(ctx context.Context, req *pb.QEMUStopRequest) (*pb.QEMUStopResponse, error) {
	// TODO: Implement QEMU stop
	return &pb.QEMUStopResponse{
		Success: true,
		Message: "QEMU stopped",
	}, nil
}

// SendInput sends input to the QEMU console.
func (s *QEMUServer) SendInput(ctx context.Context, req *pb.QEMUInputRequest) (*pb.QEMUInputResponse, error) {
	// TODO: Implement input forwarding to QEMU
	return &pb.QEMUInputResponse{
		Success: true,
	}, nil
}

// GetStatus returns the status of a QEMU instance.
func (s *QEMUServer) GetStatus(ctx context.Context, req *pb.QEMUStatusRequest) (*pb.QEMUStatus, error) {
	// TODO: Check if QEMU is running
	return &pb.QEMUStatus{
		Running:  false,
		Pid:      0,
		UptimeMs: 0,
	}, nil
}
