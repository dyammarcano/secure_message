package main

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/dyammarcano/secure_message/internal/encoding"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.0.1"

var (
	rootCmd = &cobra.Command{
		Use:   "sm",
		Short: "sm is a CLI tool to encrypt and decrypt messages",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(version)
			os.Exit(0)
		},
	}

	encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the input message",
		Long: `Encrypt command takes the message as the input.
Examples:
  sm encrypt "message"
  echo "message" | sm encrypt
  sm encrypt -i file.txt -o out.txt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Arquivo
			if inputFile := viper.GetString("input"); inputFile != "" {
				return encryptFile(cmd)
			}
			// 2. Args diretos
			if len(args) > 0 {
				encrypted, err := encoding.Serialize(argsToString(args))
				if err != nil {
					return err
				}
				cmd.Println(encrypted)
				return nil
			}
			// 3. Stdin
			stdinData, err := readStdin()
			if err != nil {
				return err
			}
			if stdinData != "" {
				encrypted, err := encoding.Serialize(stdinData)
				if err != nil {
					return err
				}
				cmd.Println(encrypted)
				return nil
			}
			return errors.New("no input provided (args, file, or stdin)")
		},
	}

	decryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypts the input message",
		Long: `Decrypt command takes the encrypted message as the input.
Examples:
  sm decrypt "encrypted message"
  echo "encrypted message" | sm decrypt
  sm decrypt -i file.enc -o file.txt`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Arquivo
			if inputFile := viper.GetString("input"); inputFile != "" {
				return decryptFile(cmd)
			}
			// 2. Args diretos
			if len(args) > 0 {
				decrypted, err := encoding.Deserialize(argsToString(args))
				if err != nil {
					return err
				}
				cmd.Println(decrypted)
				return nil
			}
			// 3. Stdin
			stdinData, err := readStdin()
			if err != nil {
				return err
			}
			if stdinData != "" {
				decrypted, err := encoding.Deserialize(stdinData)
				if err != nil {
					return err
				}
				cmd.Println(decrypted)
				return nil
			}
			return errors.New("no input provided (args, file, or stdin)")
		},
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cobra.CheckErr(rootCmd.ExecuteContext(ctx))
}

func init() {
	rootCmd.Flags().StringP("input", "i", "", "input file")
	rootCmd.Flags().StringP("output", "o", "", "output file")

	// 🔧 Bind flags to viper
	_ = viper.BindPFlag("input", rootCmd.Flags().Lookup("input"))
	_ = viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))

	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func encryptFile(cmd *cobra.Command) error {
	inputFile := viper.GetString("input")

	data, err := os.ReadFile(inputFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return err
		}

		return err
	}

	encrypted, err := encoding.Serialize(string(data))
	if err != nil {
		return err
	}

	outputFile := viper.GetString("output")
	if outputFile != "" {
		if err = os.WriteFile(outputFile, []byte(encrypted), 0644); err != nil {
			return err
		}

		return nil
	}

	cmd.Println(encrypted)

	return nil
}

func decryptFile(cmd *cobra.Command) error {
	inputFile := viper.GetString("input")

	file, err := os.ReadFile(inputFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return err
		}

		return err
	}

	decrypted, err := encoding.Deserialize(string(file))
	if err != nil {
		return err
	}

	outputFile := viper.GetString("output")
	if outputFile != "" {
		if err = os.WriteFile(outputFile, []byte(decrypted), 0644); err != nil {
			return err
		}

		return nil
	}

	cmd.Println(decrypted)

	return nil
}

// argsToString converts args to string
func argsToString(args []string) string {
	return strings.TrimSpace(strings.Join(args, " "))
}

// readStdin reads data piped from stdin, returns "" if nothing
func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(data)), nil
	}

	return "", nil
}
