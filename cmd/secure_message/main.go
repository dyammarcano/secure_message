package main

import (
	"context"
	"fmt"
	"github.com/dyammarcano/secure_message/internal/metadata"
	"github.com/spf13/cobra"
	"os"
	"runtime/trace"
)

var (
	Version    = "v0.0.1-manual-build"
	CommitHash string
	Date       string
)

var rootCmd = &cobra.Command{
	Use:   "secure_message",
	Short: "secure_message is a CLI tool to encrypt and decrypt messages",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, _ []string) {
		defer trace.StartRegion(cmd.Context(), "version").End()
		// clean console
		fmt.Fprintf(cmd.OutOrStdout(), "\033[H\033[2J")

		// print version
		fmt.Fprintf(cmd.OutOrStdout(), metadata.String())
		os.Exit(0)
	},
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rootCmd.ExecuteContext(ctx)
	cobra.CheckErr(err)
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	metadata.Set(Version, CommitHash, Date)
}
