package cmd

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gtoo/port"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var host string
var ports string
var timeout int
var threads int
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
		var portScanner port.PortScanner
		err := portScanner.SetHosts(args[0])
		if err != nil {
			panic(err)
		}
		log.Printf("set host %s.\n", args[0])
		portScanner.SetThreads(threads)
		log.Printf("set threads %d.\n", threads)
		portScanner.SetTimeout(time.Duration(timeout) * time.Second)
		log.Printf("set timeout %d s.\n", timeout)
		err = portScanner.SetPort(ports)
		if err != nil {
			panic(err)
		}

		opened := portScanner.GetOpenedPort()
		log.Printf("result: %v.\n", strings.Replace(strings.Trim(fmt.Sprint(opened), "[]"), " ", ",", -1))
	},
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().StringVarP(&ports, "port", "p", "1-65535", "scan port range")
	portCmd.Flags().IntVarP(&threads, "threads", "t", 1000, "scan threads")
	portCmd.Flags().IntVarP(&timeout, "timeout", "", 3, "scan timeout")
}
