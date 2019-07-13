package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/benmatselby/frost/version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// NewRootCommand will return the application
func NewRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "knope",
		Short:   "Get information out of AWS CodeBuild",
		Version: version.GITCOMMIT,
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.benmatselby/knope.yaml)")

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 viper.GetString("AWS_PROFILE"),
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	})
	if err != nil {
		os.Exit(2)
	}

	// Create CodeBuild service client
	svc := codebuild.New(sess)

	cmd.AddCommand(
		NewListProjectsCommand(svc),
	)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	initConfig()

	cmd := NewRootCommand()

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
