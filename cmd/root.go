package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"gtoo/convert"
	"gtoo/domain"
	"gtoo/ip"
	"gtoo/utils"
	"io"
	"net"
	"os"
	"path"
	"strings"

	"github.com/projectdiscovery/cdncheck"
	"github.com/saucer-man/iplookup"
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
		log.Info("微步在线查询IP信息:")
		if ThreadBookAPIKey != "" {
			t := ip.NewThreatBook(ThreadBookAPIKey)
			err := t.IpReputation(address.String())
			if err != nil {
				log.Errorf("微步在线查询IP信息失败: %v", err)
			}
		} else {
			log.Warn("未提供ThreadBookAPIKey")
		}
		// cdn查询
		// uses projectdiscovery endpoint with cached data to avoid ip ban
		// Use cdncheck.New() if you want to scrape each endpoint (don't do it too often or your ip can be blocked)
		log.Info("cdn查询(不一定准):")
		client, err := cdncheck.NewWithCache()
		if err != nil {
			log.Warnf("查询发生错误:%v", err)
		}
		found, err := client.Check(address)
		if err != nil {
			log.Warnf("查询发生错误:%v", err)
		} else {
			if found {
				log.Info("此ip属于cdn网段")
			} else {
				log.Info("此ip不属于cdn网段")
			}
		}

		log.Info("ip138查询ip信息:")
		err = ip.Ipinfo(address.String())
		if err != nil {
			log.Warnf("查询发生错误:%v", err)
		}
		log.Info("ip反查域名:")
		err = ip.IpLookup(address.String())
		if err != nil {
			log.Warnf("查询发生错误:%v", err)
		}
		log.Info("ip反查域名版本2:")
		err = iplookup.LookUp(address.String())
		if err != nil {
			log.Warnf("查询发生错误:%v", err)
		}
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
		log.Infof("ICP备案结果:")
		ipc138, err := domain.IcpChinaz(d)
		if err != nil {
			log.Debugf("Ip138没查到备案信息：%v", err)
			ipcvvhan, err := domain.IcpVvhan(d)
			if err != nil {
				log.Info(err)
			} else {
				utils.PrintUseTag(ipcvvhan)
			}

		} else {

			utils.PrintUseTag(ipc138)

		}
		log.Info("企业工商信息：")
		entinfo, err := domain.EntInfo(d)
		if err != nil {
			log.Warnf("企业信息查询出错: %v", err)
		} else {
			utils.PrintUseTag(entinfo)
		}
		log.Info("企业域名查询")
		domains, err := domain.GetCompanyDomain(entinfo.CompanyName)
		if err != nil {
			log.Warnf("企业域名查询: %v", err)
		} else {
			fmt.Printf("一共%d个域名:\n", len(domains))
			for _, domain := range domains {
				fmt.Println(domain.Host)
			}

		}
	},
}
var output string
var d string //domain
var domainFile string
var subdomainCmd = &cobra.Command{
	Use: "subdomain",
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	// if len(args) < 1 {
	// 	// 	return errors.New("requires at least args\nExample: gtoo domain subdomain example.com")
	// 	// }
	// 	return nil
	// },
	Run: func(cmd *cobra.Command, args []string) {
		var domains []string
		if d != "" {
			d = utils.GetDomain(d)
			if d != "" {
				domains = append(domains, d)
			}

		}
		if domainFile != "" {
			fi, err := os.Open(domainFile)
			if err != nil {
				log.Fatal("读取文件失败：%s", domainFile)
			}
			defer fi.Close()
			br := bufio.NewReader(fi)
			for {
				line, _, err := br.ReadLine()
				if err == io.EOF {
					break
				}
				l := utils.GetDomain(string(line))
				if l != "" {
					domains = append(domains, l)
				}
			}
		}
		if len(domains) == 0 {
			log.Fatal("请使用--domain或者--domain-file指定扫描目标")
		}
		// 设置一下输出文件
		dir, _ := os.Getwd()
		outputPath := path.Join(dir, output)
		outputDir := path.Dir(outputPath)
		if !utils.Exists(outputDir) {
			err := os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				log.Fatalf("设置输出路径失败: %v", err)
			}
		}
		log.Infof("设置输出文件: %s", outputPath)

		// 设置一下加载的api
		apistr := ""
		for k := range domain.SubdomainApi {
			apistr = fmt.Sprintf("%s [%s]", apistr, k)
		}

		// 开始扫描
		log.Infof("下面开始api扫描:%s", apistr)
		domain.ApiScan(domains, output)
		log.Infof("api扫描结束!")
	},
}

var rootCmd = &cobra.Command{
	Use:   "gtoo",
	Short: "gtoo id a pentese tool",
	Long:  `gtoo id a pentese tool, when you have this, you will got all`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}
	},
}
var Verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(convertCmd)
	convertCmd.AddCommand(base64encodeCmd)
	convertCmd.AddCommand(base64decodeCmd)
	convertCmd.AddCommand(md5Cmd)
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(domaininfoCmd)
	domainCmd.AddCommand(subdomainCmd)
	subdomainCmd.Flags().StringVarP(&output, "output", "", "result.txt", "output file")
	subdomainCmd.Flags().StringVarP(&d, "domain", "d", "", "domain to scan")
	subdomainCmd.Flags().StringVarP(&domainFile, "domain-file", "", "", "domain file to scan")

	rootCmd.AddCommand(versionCmd)

	ipInfoCmd.Flags().StringVarP(&ThreadBookAPIKey, "threadbookapikey", "", "", "微步在线apikey")
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
