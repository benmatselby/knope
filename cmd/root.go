package cmd

import (
	"fmt"
	"os"

	"github.com/benmatselby/knope/client"
	"github.com/benmatselby/knope/version"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// NewRootCommand will return the application
func NewRootCommand(client client.API) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "knope",
		Short:   "CLI tool for retrieving data from AWS CodeBuild",
		Version: version.GITCOMMIT,
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.benmatselby/knope.yaml)")

	cmd.AddCommand(
		NewListProjectsCommand(client),
		NewListBuildsForProjectCommand(client),
	)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	initConfig()

	client := client.NewClient()

	cmd := NewRootCommand(&client)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

		viper.AddConfigPath(home + ".benmatselby/")
		viper.SetConfigName("knope")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
