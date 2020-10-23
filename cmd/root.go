package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gtoo",
	Short: "gtoo id a pentese tool",
	Long:  `gtoo id a pentese tool, when you have this, you will got all`,
}

func Execute() {
	rootCmd.Execute()
}
