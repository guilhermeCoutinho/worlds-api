package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start metagame",
	Long:  `start metagame`,
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

func StartServer() {
	fmt.Println("Starting server...")
}

func init() {
	rootCmd.AddCommand(startCmd)
}
