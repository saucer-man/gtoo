package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"gtoo/port"
	"gtoo/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ports string
var masscanRate int
var masscanPath string
var includeFile string
var isProbe bool
var scanProbeFile string
var routines int
var scanRarity int
var timeout int

const nomalPort = "1,7,9,13,19,21-23,25,37,42,49,53,69,79-81,85,105,109-111,113,123,135,137-139,143,161,179,222,264,384,389,402,407,443-446,465,500,502,512-515,523-524,540,548,554,587,617,623,689,705,771,783,873,888,902,910,912,921,993,995,998,1000,1024,1030,1035,1090,1098-1103,1128-1129,1158,1199,1211,1220,1234,1241,1300,1311,1352,1433-1435,1440,1494,1521,1530,1533,1581-1582,1604,1720,1723,1755,1811,1900,2000-2001,2049,2082,2083,2100,2103,2121,2199,2207,2222,2323,2362,2375,2380-2381,2525,2533,2598,2601,2604,2638,2809,2947,2967,3000,3037,3050,3057,3128,3200,3217,3273,3299,3306,3311,3312,3389,3460,3500,3628,3632,3690,3780,3790,3817,4000,4322,4433,4444-4445,4659,4679,4848,5000,5038,5040,5051,5060-5061,5093,5168,5247,5250,5351,5353,5355,5400,5405,5432-5433,5498,5520-5521,5554-5555,5560,5580,5601,5631-5632,5666,5800,5814,5900-5910,5920,5984-5986,6000,6050,6060,6070,6080,6082,6101,6106,6112,6262,6379,6405,6502-6504,6542,6660-6661,6667,6905,6988,7001,7021,7071,7080,7144,7181,7210,7443,7510,7579-7580,7700,7770,7777-7778,7787,7800-7801,7879,7902,8000-8001,8008,8014,8020,8023,8028,8030,8080-8082,8087,8090,8095,8161,8180,8205,8222,8300,8303,8333,8400,8443-8444,8503,8800,8812,8834,8880,8888-8890,8899,8901-8903,9000,9002,9060,9080-9081,9084,9090,9099-9100,9111,9152,9200,9390-9391,9443,9495,9809-9815,9855,9999-10001,10008,10050-10051,10080,10098,10162,10202-10203,10443,10616,10628,11000,11099,11211,11234,11333,12174,12203,12221,12345,12397,12401,13364,13500,13838,14330,15200,16102,17185,17200,18881,19300,19810,20010,20031,20034,20101,20111,20171,20222,22222,23472,23791,23943,25000,25025,26000,26122,27000,27017,27888,28222,28784,30000,30718,31001,31099,32764,32913,34205,34443,37718,38080,38292,40007,41025,41080,41523-41524,44334,44818,45230,46823-46824,47001-47002,48899,49152,50000-50004,50013,50500-50504,52302,55553,57772,62078,62514,65535"

var portCmd = &cobra.Command{
	Use:   "portscan",
	Short: "scan open port of host",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one host to scan")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		m := port.New()
		// 判断是否为root
		if !utils.IsRoot() {
			log.Warnf("need to sudo or run as root or something")
			os.Exit(-1)
		}
		// 设置masscan路径
		path, err := exec.LookPath(masscanPath)
		if err != nil {
			log.Warnf("didn't find '{%s}' executable", masscanPath)
			os.Exit(-1)
		}
		m.SetSystemPath(path)

		// 扫描目标范围
		m.SetTargets(args[0])
		if includeFile != "" {
			if !utils.Exists(includeFile) {
				log.Warnf("didn't find '{%s}' target file", includeFile)
				os.Exit(-1)
			}
			m.SetTargetFile(includeFile)
		}
		// 设置扫描端口范围
		if ports == "" {
			ports = nomalPort
		}
		m.SetPorts(ports)

		// 扫描速率
		m.SetRate(masscanRate)

		// 设置输出路径
		rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		outputDir := filepath.Join(rootDir, "output")
		if !utils.IsDir(outputDir) {
			os.Mkdir(outputDir, os.ModePerm)
		}
		masscanOutput := filepath.Join(outputDir, "masscan_result.xml")
		m.SetOutput(masscanOutput)
		log.Info("port scan with masscan ...")
		// 开始扫描
		err = m.Run()
		if err != nil {
			log.Warnf("scanner failed: %v", err)
			os.Exit(-1)
		}

		// 解析扫描结果
		results, err := m.Parse(masscanOutput)
		if err != nil {
			log.Warnf("Parse result failed: %v", err)
			os.Exit(-1)
		}
		maascanRes := filepath.Join(outputDir, "masscan_result.txt")
		f, _ := os.Create(maascanRes) //创建文件
		defer f.Close()
		maascanResFile := bufio.NewWriter(f)
		for _, result := range results {
			for _, port := range result.Ports {
				maascanResFile.WriteString(fmt.Sprintf("%s:%s\n", result.Address.Addr, port.Portid))
			}
		}
		maascanResFile.Flush()
		log.Infof("maascan result save at %s", maascanRes)

		// 如果不进行探测则直接退出
		if !isProbe {
			os.Exit(0)
		}

		// 进行指纹探测
		log.Info("probe detecting...")
		probeRes := filepath.Join(outputDir, "probe_result.txt")
		probef, _ := os.Create(probeRes) //创建文件
		defer probef.Close()
		probeF := bufio.NewWriter(probef)
		v := port.VScan{}
		v.Init(scanProbeFile)
		// 输入输出缓冲为最大协程数量的 5 倍
		inTargetChan := make(chan port.Target, routines*5)
		outResultChan := make(chan port.Result, routines*2)
		// 最大协程并发量为参数 routines
		wgWorkers := sync.WaitGroup{}
		wgWorkers.Add(int(routines))
		// config
		var config port.Config
		config.Rarity = scanRarity
		config.SendTimeout = time.Duration(timeout) * time.Second
		config.ReadTimeout = time.Duration(timeout) * time.Second

		config.UseAllProbes = false
		config.NULLProbeOnly = false
		// 启动协程并开始监听处理输入的 Target
		for i := 0; i < routines; i++ {
			worker := port.Worker{
				In:     inTargetChan,
				Out:    outResultChan,
				Config: &config,
			}
			worker.Start(&v, &wgWorkers)
		}
		// 实时结果输出协程
		wgOutput := sync.WaitGroup{}
		wgOutput.Add(1)
		go func(wg *sync.WaitGroup) {
			probeF.WriteString("ip\tport\tservice\tversion\tproduct\tbanner\n")
			for {
				result, ok := <-outResultChan
				if ok {
					service := result.Service.Name
					// extras := result.Service.Extras
					version := result.Service.Info
					product := result.Service.VendorProduct
					banner := result.Service.Banner
					banner = strings.Replace(banner, "\r", "", -1)
					banner = strings.Replace(banner, "\n", "", -1)
					// version := fmt.Sprintf("%s%s%s%s%s%s%s", extras.VendorProduct, extras.Version, extras.Info, extras.Hostname, extras.OperatingSystem, extras.DeviceType, extras.CPE)
					probeF.WriteString(fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\n", result.IP, result.Port, service, version, product, banner))
					// probeF.Write([]byte(banner))
					// encodeJSON, err := json.Marshal(result)
					// if err != nil {
					// 	continue
					// }
					// probeF.WriteString(string(encodeJSON) + "\n")
				} else {
					break
				}
			}
			probeF.Flush()
			wg.Done()
		}(&wgOutput)
		for _, result := range results {
			for _, p := range result.Ports {
				intPort, _ := strconv.Atoi(p.Portid)
				target := port.Target{
					IP:       result.Address.Addr,
					Port:     intPort,
					Protocol: p.Protocol,
				}
				inTargetChan <- target
			}
		}
		close(inTargetChan)
		wgWorkers.Wait()
		close(outResultChan)
		wgOutput.Wait()
		log.Infof("probe result save at %s", probeRes)
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().StringVarP(&masscanPath, "masscan-path", "", "masscan", "masscan path, default masscan")
	portCmd.Flags().StringVarP(&includeFile, "include-file", "f", "", "load target  from file")
	portCmd.Flags().StringVarP(&ports, "port", "p", "", "scan port range, default 1-65535 ")
	portCmd.Flags().IntVarP(&masscanRate, "masscan-rate", "", 1000, "masscan rate, default 1000")
	portCmd.Flags().BoolVarP(&isProbe, "is-probe", "", true, "is Probe? default true")

	portCmd.Flags().StringVarP(&scanProbeFile, "scan-probe-file", "", "./source/nmap-service-probes", "scan port range")
	portCmd.Flags().IntVarP(&routines, "routines", "", 10, "Goroutines numbers using during probe scanning")
	portCmd.Flags().IntVarP(&scanRarity, "scan-rarity", "", 7, "Sets the intensity level of a version scan to the specified value (default 7)")
	portCmd.Flags().IntVarP(&timeout, "timeout", "", 5, "connection timeout in seconds (default 5)")
}
