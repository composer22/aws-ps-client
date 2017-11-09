package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Used globally for all commands
var (
	cfgFile         string
	awsAccessKey    string
	awsAccessSecret string
	awsRegion       string
	format          string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "aws-ps-client",
	Short: "Retrieve values from AWS EC2 Parameter Store",
	Long:  "Client for retrieving values from AWS Parameter store for a given key",
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-ps-client.yaml)")
	RootCmd.PersistentFlags().StringP("aws-access-key", "k", "", "AWS IAM access authentication key or full path to a file containing key")
	RootCmd.PersistentFlags().StringP("aws-access-secret", "s", "", "AWS IAM access authentication secret or full path to a file containing secret")
	RootCmd.PersistentFlags().StringP("aws-region", "r", "us-west-2", "AWS region ex: us-west-2")
	RootCmd.PersistentFlags().StringP("format", "f", "bash", "Return 'bash', 'json', or 'text'")

	// Get values from config file.
	viper.BindPFlag("aws-access-key", RootCmd.PersistentFlags().Lookup("aws-access-key"))
	viper.BindPFlag("aws-access-secret", RootCmd.PersistentFlags().Lookup("aws-access-secret"))
	viper.BindPFlag("aws-region", RootCmd.PersistentFlags().Lookup("aws-region"))
	viper.BindPFlag("format", RootCmd.PersistentFlags().Lookup("format"))
	viper.SetDefault("aws-region", "us-west-2")
	viper.SetDefault("format", "text")
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
		viper.AddConfigPath(home) // adding home directory as first search path.
		viper.AddConfigPath(".")
		viper.SetConfigName(".aws-ps-client") // name of config file (without extension).
	}

	viper.AutomaticEnv() // read in environment variables that match.

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Cannot find configuration file.\n\nERR: %s\n", err.Error())
		os.Exit(0)
	}
	awsAccessKey = viper.GetString("aws-access-key")
	awsAccessSecret = viper.GetString("aws-access-secret")
	awsRegion = viper.GetString("aws-region")
	format = viper.GetString("format")

	// We allow for some config settings via a filepath, so that it loads from
	// Docker secrets mechanism.

	// Try to load awsAccessKey from a file. If successful, replace value.
	if v, err := ioutil.ReadFile(awsAccessKey); err == nil {
		awsAccessKey = string(v)
	}
	// Try to load awsAccessSecret from a file. If successful, replace value.
	if v, err := ioutil.ReadFile(awsAccessSecret); err == nil {
		awsAccessSecret = string(v)
	}
	// Try to load awsRegion from a file. If successful, replace value.
	if v, err := ioutil.ReadFile(awsRegion); err == nil {
		awsRegion = string(v)
	}
}
