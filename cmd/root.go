package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dobefu/gofutz/internal/server"
	"github.com/Dobefu/gofutz/internal/websocket"
	"github.com/spf13/cobra"
)

var port int

var rootCmd = &cobra.Command{ //nolint:exhaustruct
	Use:   "gofutz",
	Short: "Run tests incrementally and interactively",
	Run:   runRootCmd,
}

// Execute executes the root command.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(
		&port,
		"port",
		"p",
		7357,
		"The port to run the test server on",
	)
}

func runRootCmd(_ *cobra.Command, _ []string) {
	srv := server.NewServer("127.0.0.1", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- srv.Start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			slog.Error(fmt.Sprintf("Could not start server: %s", err.Error()))
			os.Exit(1)
		}
	case sig := <-sigChan:
		slog.Info(fmt.Sprintf("Received shutdown signal: %s", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		websocket.CloseAll()
		err := srv.Shutdown(ctx)

		if err != nil {
			slog.Error(fmt.Sprintf("Server shutdown error: %s", err.Error()))
			defer os.Exit(1)
		}
	}
}
