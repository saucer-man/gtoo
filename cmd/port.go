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
		log.Infof("masscan result save at {%s}", masscanOutput)
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
		finalRes := filepath.Join(outputDir, "port-result.txt")
		f, _ := os.Create(finalRes) //创建文件
		defer f.Close()
		outFile := bufio.NewWriter(f)

		// 如果不进行探测则将结果保存一下
		if !isProbe {
			for _, result := range results {
				for _, port := range result.Ports {
					outFile.WriteString(fmt.Sprintf("%s:%s\n", result.Address.Addr, port.Portid))
				}
			}
			outFile.Flush()
			log.Infof("final result save at %s", finalRes)
			os.Exit(0)
		}

		// 进行指纹探测
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
			for {
				result, ok := <-outResultChan
				if ok {
					service := result.Service.Name
					extras := result.Service.Extras
					version := fmt.Sprintf("%s%s%s%s%s%s%s", extras.VendorProduct, extras.Version, extras.Info, extras.Hostname, extras.OperatingSystem, extras.DeviceType, extras.CPE)
					outFile.WriteString(fmt.Sprintf("%s %d %s %s\n", result.IP, result.Port, service, version))
				} else {
					break
				}
			}
			outFile.Flush()
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
		log.Infof("final result save at %s", finalRes)
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().StringVarP(&masscanPath, "masscan-path", "m", "masscan", "scan port range")
	portCmd.Flags().StringVarP(&includeFile, "include-file", "f", "", "scan target filename")
	portCmd.Flags().StringVarP(&ports, "port", "p", "1-65535", "scan port range")
	portCmd.Flags().IntVarP(&masscanRate, "masscan-rate", "", 1000, "masscan rate, default 1000")
	portCmd.Flags().BoolVarP(&isProbe, "is-probe", "", true, "is Probe? default true")

	portCmd.Flags().StringVarP(&scanProbeFile, "scan-probe-file", "", "./source/nmap-service-probes", "scan port range")
	portCmd.Flags().IntVarP(&routines, "routines", "", 10, "Goroutines numbers using during probe scanning")
	portCmd.Flags().IntVarP(&scanRarity, "scan-rarity", "", 7, "Sets the intensity level of a version scan to the specified value (default 7)")
	portCmd.Flags().IntVarP(&timeout, "timeout", "", 5, "connection timeout in seconds (default 5)")
}
