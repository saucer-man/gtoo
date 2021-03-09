package cmd

import (
	"errors"
	"fmt"
	"gtoo/collect/zoomeye"
	"gtoo/config"
	"os"

	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "collect domain and ip",
}

var zoomeyeCmd = &cobra.Command{
	Use: "zoomeye",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args\nExample: gtoo collect zoomeye search_content")
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
		err = zoomeye.Search(conf.Zoomeye.APIKey, args[0])
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(collectCmd)
	collectCmd.AddCommand(zoomeyeCmd)

}
