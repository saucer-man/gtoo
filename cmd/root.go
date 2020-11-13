package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gtoo",
	Short: "gtoo id a pentese tool",
	Long:  `gtoo id a pentese tool, when you have this, you will got all`,
}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of gtoo",
	Long:  `version of gtoo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gtoo version: v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	rootCmd.Execute()
}
