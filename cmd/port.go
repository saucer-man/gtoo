package cmd

import (
	"errors"
	"gtoo/port"
	"gtoo/utils"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var host string
var ports string
var timeout int
var threads int
var masscanPath string
var includeFile string
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
		m.SetRate(threads)

		// 开始扫描
		err = m.Run()
		if err != nil {
			log.Warnf("scanner failed: %v", err)
			os.Exit(-1)
		}

		// 解析扫描结果
		results, err := m.Parse()
		if err != nil {
			log.Warnf("Parse result failed: %v", err)
			os.Exit(-1)
		}

		for _, result := range results {
			log.Info(result)
		}
		// var portScanner port.PortScanner
		// err := portScanner.SetHosts(args[0])
		// if err != nil {
		// 	panic(err)
		// }
		// log.Printf("set host %s.\n", args[0])
		// portScanner.SetThreads(threads)
		// log.Printf("set threads %d.\n", threads)
		// portScanner.SetTimeout(time.Duration(timeout) * time.Second)
		// log.Printf("set timeout %d s.\n", timeout)
		// err = portScanner.SetPort(ports)
		// if err != nil {
		// 	panic(err)
		// }

		// opened := portScanner.GetOpenedPort()
		// log.Printf("result: %v.\n", strings.Replace(strings.Trim(fmt.Sprint(opened), "[]"), " ", ",", -1))
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().StringVarP(&masscanPath, "masscan-path", "m", "masscan", "scan port range")
	portCmd.Flags().StringVarP(&includeFile, "includefile", "f", "", "scan target filename")
	portCmd.Flags().StringVarP(&ports, "port", "p", "1-65535", "scan port range")
	portCmd.Flags().IntVarP(&threads, "threads", "t", 1000, "scan threads or rate")
	portCmd.Flags().IntVarP(&timeout, "timeout", "", 3, "scan timeout")
}
