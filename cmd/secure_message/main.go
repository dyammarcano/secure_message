package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dyammarcano/secure_message/internal/config"
	"github.com/dyammarcano/secure_message/internal/encoding"
	"github.com/dyammarcano/secure_message/internal/helpers"
	"github.com/dyammarcano/secure_message/internal/management"
	"github.com/dyammarcano/secure_message/internal/metadata"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"reflect"
	"runtime/trace"
	"strings"
	"syscall"
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

	keysCmd = &cobra.Command{
		Use:   "keys",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if config.C.Flags.ExportFilePathKey != "" {
				exportKey(cmd, args)
			}

			if config.C.Flags.ImportFilePathKey != "" {
				importKey(cmd, args)
			}
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

	keysCmd.Flags().StringVarP(&config.C.Flags.ExportFilePathKey, "export", "e", "", "export key to file")
	keysCmd.Flags().StringVarP(&config.C.Flags.ImportFilePathKey, "import", "i", "", "import key from file")

	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(keysCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	metadata.Set(Version, CommitHash, Date)
}

func getPassword(cmd *cobra.Command, _ []string, confirm bool) ([]byte, error) {
	cmd.Print("Enter password:")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	if confirm {
		cmd.Print("\nConfirm password:")
		bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}

		if !reflect.DeepEqual(bytePassword, bytePassword) {
			return nil, errors.New("passwords do not match")
		}
	}

	return bytePassword, nil
}

func exportKey(cmd *cobra.Command, args []string) {
	password, err := getPassword(cmd, args, true)
	if err != nil {
		printAndExit(cmd, args, err.Error())
	}

	data, err := management.ExportKeys()
	if err != nil {
		printAndExit(cmd, args, err.Error())
	}

	bytesData, err := encoding.SerializeWithPassword(data, password)
	if err != nil {
		printAndExit(cmd, args, err.Error())
	}

	if err := os.WriteFile(config.C.Flags.ExportFilePathKey, bytesData, 0644); err != nil {
		printAndExit(cmd, args, err.Error())
	}

	cmd.Println("Keys exported successfully")
}

func importKey(cmd *cobra.Command, args []string) {
	if err := warningMessage(cmd, args); err != nil {
		printAndExit(cmd, args, err.Error())
	}

	password, err := getPassword(cmd, args, false)
	if err != nil {
		printAndExit(cmd, args, err.Error())
	}

	data, err := os.ReadFile(config.C.Flags.ImportFilePathKey)
	if err != nil {
		printAndExit(cmd, args, err.Error())
	}

	data, err = encoding.DeserializeWithPassword(data, password)
	if err != nil {
		printAndExit(cmd, args, err.Error())
	}

	if err = management.ImportKeys(data); err != nil {
		printAndExit(cmd, args, err.Error())
	}

	cmd.Println("Keys imported successfully")
}

func warningMessage(cmd *cobra.Command, _ []string) error {
	cmd.Print("Warning: This command will overwrite your current keys\nDo you want to continue? [y/N] ")

	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		return err
	}

	if strings.ToLower(answer) != "y" {
		return errors.New("aborted")
	}

	return nil
}

func printAndExit(cmd *cobra.Command, _ []string, message string) {
	cmd.Println(message)
	os.Exit(0)
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
