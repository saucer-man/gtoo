package cmd

import (
	"errors"
	"gtoo/domain"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "search some domain info",
}
var infoCmd = &cobra.Command{
	Use: "info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo domain info example.com")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		d := args[0]
		if !strings.HasPrefix(d, "http") {
			d = "https://" + d
		}
		err := domain.Whois(d)
		if err != nil {
			log.Errorf("whois查询出错: %v", err)
		}
		err = domain.Ipc(d)
		if err != nil {
			log.Errorf("IPC备案查询出错: %v", err)
		}
		// TODO is cdn
		// TODO 威胁情报
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(infoCmd)
}
