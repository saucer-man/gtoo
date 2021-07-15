package cmd

import (
	"errors"
	"fmt"
	"gtoo/convert"
	"gtoo/domain"
	"gtoo/ip"
	"gtoo/utils"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "ip相关信息",
}
var ThreadBookAPIKey string
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

		if ThreadBookAPIKey != "" {
			t := ip.NewThreatBook(ThreadBookAPIKey)
			log.Info("微步在线查询IP信息...")
			err := t.IP(address.String())
			if err != nil {
				log.Errorf("微步在线查询IP信息失败: %v", err)
			}
		}
		// TODO is cdn
	},
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "常用编码解码",
}

var base64encodeCmd = &cobra.Command{
	Use: "base64encode",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo convert base64encode str_to_encode")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		str := convert.Base64encode([]byte(args[0]))
		fmt.Println(str)
	},
}
var base64decodeCmd = &cobra.Command{
	Use: "base64decode",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo convert base64encode str_to_decode")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		data, err := convert.Base64decode(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	},
}

var md5Cmd = &cobra.Command{
	Use: "md5encode",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo convert base64encode str_to_encode")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if utils.Exists(args[0]) {
			data, err := convert.Md5encodeFile(args[0])
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s file's md5 is %s\n", args[0], data)
		} else {
			data, err := convert.Md5encodeStr(args[0])
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s str's md5 is %s\n", args[0], data)
		}
	},
}

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "域名相关信息",
}
var domaininfoCmd = &cobra.Command{
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
			d = "http://" + d
		}
		u, err := utils.ParseURL(d)
		if err != nil {
			log.Errorf("解析域名出错: %v", err)
		}
		d = u.Domain + "." + u.TLD
		log.Infof("解析出域名:%s\n", d)
		err = domain.Whois(d)
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
var subdomainCmd = &cobra.Command{
	Use: "subdomain",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo domain subdomain example.com")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		d := args[0]
		if !strings.HasPrefix(d, "http") {
			d = "http://" + d
		}
		u, err := utils.ParseURL(d)
		if err != nil {
			log.Errorf("解析域名出错: %v", err)
		}
		d = u.Domain + "." + u.TLD
		log.Infof("解析出域名:%s", d)
		log.Infof("%s 是否泛解析: %t", d, domain.IsWildCard(d))
		subdomains, err := domain.MatchSubdomains(d, "bbb.aaa.seebug.com www.seebug.com")
		if err != nil {
			log.Errorf("解析域名出错: %v", err)
		}
		log.Info(subdomains)
		// err := domain.Whois(d)
		// if err != nil {
		// 	log.Errorf("whois查询出错: %v", err)
		// }
		// err = domain.Ipc(d)
		// if err != nil {
		// 	log.Errorf("IPC备案查询出错: %v", err)
		// }
		// TODO is cdn
		// TODO 威胁情报
	},
}

var rootCmd = &cobra.Command{
	Use:   "gtoo",
	Short: "gtoo id a pentese tool",
	Long:  `gtoo id a pentese tool, when you have this, you will got all`,
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.AddCommand(base64encodeCmd)
	convertCmd.AddCommand(base64decodeCmd)
	convertCmd.AddCommand(md5Cmd)

	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(domaininfoCmd)
	domainCmd.AddCommand(subdomainCmd)

	rootCmd.AddCommand(versionCmd)

	ipInfoCmd.PersistentFlags().StringVarP(&ThreadBookAPIKey, "threadbookapikey", "", "", "微步在线apikey")
	rootCmd.AddCommand(ipCmd)
	ipCmd.AddCommand(ipInfoCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of gtoo",
	Long:  `version of gtoo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gtoo version: v0.1")
	},
}

func Execute() {
	rootCmd.Execute()
}
