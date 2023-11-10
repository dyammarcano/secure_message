package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dyammarcano/secure_message/internal/config"
	"github.com/dyammarcano/secure_message/internal/encoding"
	"github.com/dyammarcano/secure_message/internal/helpers"
	"github.com/dyammarcano/secure_message/internal/metadata"
	"github.com/spf13/cobra"
	"os"
	"runtime/trace"
)

var (
	Version    = "v0.0.1-manual-build"
	CommitHash string
	Date       string

	rootCmd = &cobra.Command{
		Use:   "secure_message",
		Short: "secure_message is a CLI tool to encrypt and decrypt messages",
	}

	versionCmd = &cobra.Command{
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

	encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the input message",
		Long: `Encrypt command takes the message as the input 
./my_app decrypt "encrypt message to encrypt"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if config.C.Flags.InputFile != "" {
				return encryptFile(cmd, args)
			}

			msg := helpers.ArgsToString(args)

			// encrypt message
			encrypted, err := encoding.Serialize(msg)
			if err != nil {
				return err
			}

			// print encrypted message
			cmd.Println(encrypted)

			return nil
		},
	}

	decryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypts the input message",
		Long: `Decrypt command takes the encrypted message as the input 
and provides the decrypted message as the output. For example:

./secure_message decrypt "encrypted message"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !helpers.FileExists(config.C.DatabaseFilePath) {
				cmd.Println("Database file not found, import keys to continue.")
				os.Exit(1)
			}

			if config.C.Flags.InputFile != "" {
				return decryptFile(cmd, args)
			}

			msg := helpers.ArgsToString(args)

			// decrypt message
			decrypted, err := encoding.Deserialize(msg)
			if err != nil {
				return err
			}

			// print decrypted message
			cmd.Println(decrypted)
			return nil
		},
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rootCmd.ExecuteContext(ctx)
	cobra.CheckErr(err)
}

func init() {
	decryptCmd.Flags().StringVarP(&config.C.Flags.InputFile, "input", "i", "", "input file to decrypt")
	decryptCmd.Flags().StringVarP(&config.C.Flags.OutputFile, "output", "o", "", "output file to save decrypted message")

	encryptCmd.Flags().StringVarP(&config.C.Flags.InputFile, "input", "i", "", "input file to encrypt")
	encryptCmd.Flags().StringVarP(&config.C.Flags.OutputFile, "output", "o", "", "output file to save encrypted message")

	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	metadata.Set(Version, CommitHash, Date)
}

func encryptFile(cmd *cobra.Command, args []string) error {
	// check if file exists
	if !helpers.FileExists(config.C.Flags.InputFile) {
		return errors.New("file not found")
	}

	// read file
	file, err := os.ReadFile(config.C.Flags.InputFile)
	if err != nil {
		return err
	}

	// encrypt file
	encrypted, err := encoding.Serialize(string(file))
	if err != nil {
		return err
	}

	if config.C.Flags.OutputFile != "" {
		// write encrypted message to file
		err = os.WriteFile(config.C.Flags.OutputFile, []byte(encrypted), 0644)
		if err != nil {
			return err
		}
		return nil
	}

	// print encrypted message
	cmd.Println(encrypted)

	return nil
}

func decryptFile(cmd *cobra.Command, args []string) error {
	// check if file exists
	if !helpers.FileExists(config.C.Flags.InputFile) {
		return errors.New("file not found")
	}

	// read file
	file, err := os.ReadFile(config.C.Flags.InputFile)
	if err != nil {
		return err
	}

	// decrypt message
	decrypted, err := encoding.Deserialize(string(file))
	if err != nil {
		return err
	}

	if config.C.Flags.OutputFile != "" {
		// write to file
		if err := os.WriteFile(config.C.Flags.OutputFile, []byte(decrypted), 0644); err != nil {
			return err
		}
		return nil
	}

	// print decrypted message
	cmd.Println(decrypted)
	return nil
}
