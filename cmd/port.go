package cmd

import (
	"errors"
	"fmt"

	. "gtoo/port"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var host string
var port string
var portCmd = &cobra.Command{
	Use:   "portscan",
	Short: "scan open port of host",

	Run: func(cmd *cobra.Command, args []string) {
		var StartPort, EndPort int
		// 解析输入ip地址
		addr, err := net.ResolveIPAddr("ip4", host)
		if err != nil {
			fmt.Println("target parse error:", err.Error())
			os.Exit(1)
		}
		fmt.Printf("tareget addr:%v-%T\n", addr, addr)
		// 解析port地址
		if port == "" {
			StartPort = 1
			EndPort = 65535
		} else {
			StartPort, EndPort, err = ScanPortParse(port)
			if err != nil {
				fmt.Println("port parse error:", err.Error())
				os.Exit(1)
			}
		}
		fmt.Printf("StartPort:%v EndPort:%v\n", StartPort, EndPort)

		PortScan(addr, StartPort, EndPort)
	},
}

func ScanPortParse(port string) (StartPort int, EndPort int, err error) {
	//初始化返回值
	StartPort, EndPort, err = 0, 0, nil
	s := strings.Split(port, "-")
	StartPort, err = strconv.Atoi(s[0])
	if err != nil {
		fmt.Printf("cant parse scan port with %v\n", port)
		return
	}
	if len(s) == 1 {
		EndPort = StartPort + 1
	} else if len(s) == 2 {
		EndPort, err = strconv.Atoi(s[1])
		if err != nil {
			fmt.Printf("cant parse scan port with %v\n", port)
			return
		}
	} else {
		err = errors.New("input port error")
	}

	//解析完了之后开看规则
	if StartPort <= 0 || StartPort >= 65535 {
		err = errors.New("The port to start scanning is out of range")
	} else if EndPort <= 1 || EndPort >= 65536 {
		err = errors.New("The port to end scanning is out of range")
	} else if EndPort <= StartPort {
		err = errors.New("the start port of the scan cannot be larger than the end port")
	}
	return
}
func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().StringVarP(&host, "host", "t", "", "target hostname") //默认值
	portCmd.MarkFlagRequired("host")
	portCmd.Flags().StringVarP(&port, "port", "p", "", "scan port range")
}
