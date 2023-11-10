package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dyammarcano/secure_message/internal/encoding"
	"github.com/dyammarcano/secure_message/internal/metadata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"runtime/trace"
	"strings"
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
			inputFile := viper.GetString("input")
			if inputFile != "" {
				return encryptFile(cmd, args)
			}

			encrypted, err := encoding.Serialize(argsToString(args))
			if err != nil {
				return err
			}

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
			inputFile := viper.GetString("input")

			if inputFile != "" {
				return decryptFile(cmd, args)
			}

			decrypted, err := encoding.Deserialize(argsToString(args))
			if err != nil {
				return err
			}

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
	rootCmd.PersistentFlags().String("input", "", "input file")
	if err := viper.BindPFlag("input", rootCmd.Flags().Lookup("input")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().String("output", "", "output file")
	if err := viper.BindPFlag("output", rootCmd.Flags().Lookup("output")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	metadata.Set(Version, CommitHash, Date)
}

func encryptFile(cmd *cobra.Command, _ []string) error {
	inputFile := viper.GetString("input")

	if _, err := os.Stat(inputFile); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cobra.CheckErr(err)
		}
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	encrypted, err := encoding.Serialize(string(data))
	if err != nil {
		return err
	}

	outputFile := viper.GetString("ouput")

	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(encrypted), 0644)
		if err != nil {
			return err
		}
		return nil
	}

	cmd.Println(encrypted)

	return nil
}

func decryptFile(cmd *cobra.Command, args []string) error {
	inputFile := viper.GetString("input")

	if _, err := os.Stat(inputFile); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cobra.CheckErr(err)
		}
	}

	file, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	decrypted, err := encoding.Deserialize(string(file))
	if err != nil {
		return err
	}

	outputFile := viper.GetString("ouput")

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(decrypted), 0644); err != nil {
			return err
		}
		return nil
	}

	cmd.Println(decrypted)
	return nil
}

// argsToString converts args to string
func argsToString(args []string) string {
	var msg string
	for _, arg := range args {
		msg += arg + " "
	}

	msg = strings.TrimRight(msg, " ")
	return msg
}
