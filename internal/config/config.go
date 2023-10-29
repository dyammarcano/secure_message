package config

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	C *Configuration
)

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	filePath := filepath.Join(home, "secureMessage.db")

	C = &Configuration{
		Home:             home,
		DatabaseFilePath: filePath,
		AppNAme:          "secure_message",
		DirBase:          "keys",
	}
}

type (
	Configuration struct {
		Home             string `mapstructure:"home" yaml:"-"`
		DatabaseFilePath string `mapstructure:"database_file_path" yaml:"databaseFilePath"`
		AppNAme          string `mapstructure:"app_name" yaml:"appName"`
		DirBase          string `mapstructure:"dir_base" yaml:"-"`
		Flags            Flags  `mapstructure:"-" yaml:"-"`
	}

	Flags struct {
		InputFile         string
		OutputFile        string
		ExportFilePathKey string
		ImportFilePathKey string
	}
)
