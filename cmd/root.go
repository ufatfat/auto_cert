package cmd

import (
	"auto_cert/cmd/printcmd"
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cert",
	Short: "auto cert",
	Long:  `auto cert`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func init() {
	rootCmd.AddCommand(printcmd.PrintCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
