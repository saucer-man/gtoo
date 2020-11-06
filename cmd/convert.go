package cmd

import (
	"errors"
	"fmt"

	"gtoo/convert.go"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Common encoding and decoding",
}

var base64encodeCmd = &cobra.Command{
	Use: "base64encode",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args")
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
			return errors.New("requires at least args")
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

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.AddCommand(base64encodeCmd)
	convertCmd.AddCommand(base64decodeCmd)
}
