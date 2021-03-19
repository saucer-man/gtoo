package cmd

import (
	"errors"
	"gtoo/config"
	"gtoo/ip"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "search some domain info",
}
var ipInfoCmd = &cobra.Command{
	Use: "info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo ip info example.com")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//检查ip地址是否正确
		// ParseIP 这个方法 可以用来检查 ip 地址是否正确，如果不正确，该方法返回 nil
		address := net.ParseIP(args[0])
		if address == nil {
			log.Error("ip地址格式不正确")
			os.Exit(-1)
		}
		// } else {
		// 	fmt.Println("正确的ip地址", address.String())
		// }
		//下面调用微步在线的api查询ip信誉
		conf, err := config.ReadConfig()
		if err != nil {
			log.Errorf("读取配置文件失败:%v", err)
			os.Exit(-1)
		}
		if conf.ThreadBook.APIKey != "" {
			t := ip.NewThreatBook(conf.ThreadBook.APIKey)
			log.Info("微步在线查询IP信息...")
			err = t.IP(address.String())
			if err != nil {
				log.Errorf("微步在线查询IP信息失败: %v", err)
			}
		}
		// TODO is cdn
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.AddCommand(ipInfoCmd)
}
