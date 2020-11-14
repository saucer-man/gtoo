package cmd

import (
	"errors"
	"fmt"

	"gtoo/convert"
	"gtoo/utils"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Common encoding and decoding",
}

var base64encodeCmd = &cobra.Command{
	Use: "base64 encode",
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
	Use: "base64 decode",
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

var md5Cmd = &cobra.Command{
	Use: "md5 encode",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if utils.Exists(args[0]) {
			data, err := convert.Md5encodeFile(args[0])
			if err != nil {
				panic(err)
			}
			fmt.Printf("file's md5 is %s\n", data)
		} else {
			data, err := convert.Md5encodeStr(args[0])
			if err != nil {
				panic(err)
			}
			fmt.Printf("str's md5 is %s\n", data)
		}

	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.AddCommand(base64encodeCmd)
	convertCmd.AddCommand(base64decodeCmd)
	convertCmd.AddCommand(md5Cmd)
}
