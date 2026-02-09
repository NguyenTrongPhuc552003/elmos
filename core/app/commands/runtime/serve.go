package runtime

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NguyenTrongPhuc552003/elmos/api/server"
	cmdtypes "github.com/NguyenTrongPhuc552003/elmos/core/app/commands/types"
	"github.com/spf13/cobra"
)

func BuildServe(ctx *cmdtypes.Context) *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start ELMOS gRPC server",
		Long: `Start the ELMOS gRPC server to provide programmatic API access.

The server can listen on:
  - Unix socket: unix:///tmp/elmos.sock (default, local only)
  - TCP: tcp://localhost:50051 (for remote access)

Example:
  elmos serve                        # Unix socket (default)
  elmos serve --address unix:///tmp/elmos.sock
  elmos serve --address tcp://localhost:50051`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create gRPC server
			srv, err := server.NewServer(
				ctx.Config,
				ctx.AppContext,
				ctx.Exec,
				ctx.FS,
				ctx.VolumeManager,
			)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			// Setup graceful shutdown
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			// Start server in goroutine
			errCh := make(chan error, 1)
			go func() {
				ctx.Printer.Info("Starting ELMOS gRPC server on %s...\n", address)
				errCh <- srv.Serve(address)
			}()

			// Wait for shutdown signal or error
			select {
			case sig := <-sigCh:
				ctx.Printer.Info("\nReceived signal %v, shutting down...\n", sig)
				srv.Stop()
				return nil
			case err := <-errCh:
				if err != nil {
					return fmt.Errorf("server error: %w", err)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVarP(&address, "address", "a", "unix:///tmp/elmos.sock", "Server listen address (unix:// or tcp://)")

	return cmd
}
