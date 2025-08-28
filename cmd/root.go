package cmd

import (
	"fmt"
	"os"

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
	rootCmd.Flags().IntVarP(&port, "port", "p", 7357, "The port to run the test server on")
}

func runRootCmd(_ *cobra.Command, _ []string) {
	_, _ = fmt.Fprintf(os.Stdout, "Port: %d\n", port)
}
