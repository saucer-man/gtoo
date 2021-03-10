package cmd

import (
	"errors"
	"fmt"
	"gtoo/collect/zoomeye"
	"gtoo/config"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "collect domain and ip",
}

var num int
var resource string
var output string
var zoomeyeCmd = &cobra.Command{
	Use: "zoomeye",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo collect zoomeye \"port:80 nginx\" --num 100 --resource host")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error when read config")
			fmt.Println(err)
			os.Exit(2)
		}
		if output == "" {
			output = fmt.Sprintf("%s_gtoo.txt", time.Now().Format("2006-01-02_15-04-05"))
		}
		err = zoomeye.Search(conf.Zoomeye.APIKey, args[0], resource, num, output)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.AddCommand(zoomeyeCmd)
	zoomeyeCmd.Flags().StringVarP(&resource, "resource", "", "host", "scan type web/host, default host")
	zoomeyeCmd.Flags().IntVarP(&num, "num", "n", 20, "scan number, default 20")
	zoomeyeCmd.Flags().StringVarP(&output, "output", "o", "", "output filename")
}
