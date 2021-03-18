package cmd

import (
	"errors"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "domain",
	Short: "search some domain info",
}
var ipInfoCmd = &cobra.Command{
	Use: "info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo domain info example.com")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//检查ip地址是否正确
		// ParseIP 这个方法 可以用来检查 ip 地址是否正确，如果不正确，该方法返回 nil
		address := net.ParseIP(args[0])
		if address == nil {
			fmt.Println("ip地址格式不正确")
		} else {
			fmt.Println("正确的ip地址", address.String())
		}

		// TODO is cdn
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.AddCommand(ipInfoCmd)
}
