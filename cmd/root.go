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
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/httpx/runner"
	"github.com/saucer-man/iplookup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "ip相关信息",
}
var output string
var d string //domain
var domainFile string

type tomlConfig struct {
	ThreadBookAPIKey string `toml:"threadbook_apikey"`
}

var config tomlConfig
var Verbose bool
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
			config.ThreadBookAPIKey = ThreadBookAPIKey
		}

		if config.ThreadBookAPIKey != "" {
			t := ip.NewThreatBook(config.ThreadBookAPIKey)
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
var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "url相关信息",
}

var urlInfoCmd = &cobra.Command{
	Use: "info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo url info example.com/domain.txt")
		}
		return nil
	},
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Parse the command line flags and read config files
		options := runner.ParseOptions()
		options.OutputCDN = true    // -cdn
		options.ExtractTitle = true // -title
		options.OutputCName = true  // -cname
		options.StatusCode = true   // -status-code
		options.OutputMethod = true // -method
		options.TLSProbe = true     // -tls-probe
		options.TechDetect = true   // -tech-detect
		options.HTTP2Probe = true   // -http2
		options.OutputIP = true     // -ip
		options.CSPProbe = true     // - csp-probe
		options.VHost = true        // -vhost
		options.Output = "url_result.txt"
		options.NoColor = true
		// 看输入是否是文件
		if utils.FileExists(args[0]) {
			options.InputFile = args[0]
		} else {
			// Create our Temp File
			tmpFile, err := ioutil.TempFile(os.TempDir(), "gtoo-")
			if err != nil {
				log.Fatal("Cannot create temporary file", err)
			}

			log.Debug("Created tmp File: " + tmpFile.Name())

			// Example writing to the file
			_, err = tmpFile.Write([]byte(args[0]))
			if err != nil {
				log.Fatal("Failed to write to temporary file", err)
			}
			options.InputFile = tmpFile.Name()
			// Remember to clean up the file afterwards
			defer os.Remove(tmpFile.Name())
		}
		httpxRunner, err := runner.New(options)
		if err != nil {
			log.Fatal("Could not create runner: %s\n", err)
		}
		// Setup graceful exits
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				log.Info("CTRL+C pressed: Exiting\n")
				httpxRunner.Close()
				if options.ShouldSaveResume() {
					log.Info("Creating resume file: %s\n", runner.DefaultResumeFile)
					err := httpxRunner.SaveResumeConfig()
					if err != nil {
						log.Error("Couldn't create resume file: %s\n", err)
					}
				}
				os.Exit(1)
			}
		}()

		httpxRunner.RunEnumeration()
		httpxRunner.Close()
		log.Info("结果保存在url_result.txt中")
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
		filePath, _ := os.UserHomeDir()
		filePath = path.Join(filePath, ".gtoo.toml")
		toml.DecodeFile(filePath, &config)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		filePath, _ := os.UserHomeDir()
		filePath = path.Join(filePath, ".gtoo.toml")
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("failed to create/open the file")
			fmt.Println(err)
			return
		}
		if err := toml.NewEncoder(f).Encode(config); err != nil {
			// failed to encode
			return
		}
		if err := f.Close(); err != nil {
			// failed to close the file
			return
		}
	},
}

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

	urlCmd.AddCommand(urlInfoCmd)
	rootCmd.AddCommand(urlCmd)
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
