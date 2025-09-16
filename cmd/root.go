/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	logJSON bool
	Verbose int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "worlds-api",
	Short: "Worlds API",
	Long: `Worlds API is a RESTful API for the Worlds game.
It provides endpoints for managing worlds and their content.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.worlds-api.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(
		&logJSON, "log-json", "j",
		false, "log-json output mode")

	rootCmd.PersistentFlags().IntVarP(
		&Verbose, "verbose", "v", 2,
		"Verbosity level => v0: Error, v1=Warning, v2=Info, v3=Debug, v4=Trace",
	)
}
