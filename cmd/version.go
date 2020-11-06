package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of gtoo",
	Long:  `version of gtoo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gtoo version: v0.1")
	},
}

// 被rootCmd调用
func init() {
	rootCmd.AddCommand(versionCmd)
}
