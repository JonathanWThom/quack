package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/jonathanwthom/quack/storage"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type Store interface {
	Create(string) error
	Read() ([]storage.Entry, error)
	ReadByKey(string) (storage.Entry, error)
	Delete(string) error
	Update(storage.Entry) error
}

var store Store

var rootCmd = &cobra.Command{
	Use:   "quack",
	Short: "A mini journal, just for you.",
	Long: `Write journal entries, 280 characters at a time.
Entries are encrypted and can be stored locally or in the cloud.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	store = new(storage.Storage)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".quack" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".quack")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
