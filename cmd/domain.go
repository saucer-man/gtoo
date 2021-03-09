package cmd

import (
	"errors"
	"fmt"
	"gtoo/domain"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "search some domain info",
}
var whoisCmd = &cobra.Command{
	Use: "whois",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo domain whois example.com")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := domain.Whois(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(whoisCmd)
}
