package printcmd

import (
	"auto_cert/timer"
	"github.com/spf13/cobra"
)

var PrintCmd = &cobra.Command{
	Use:   "print",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		timer.Timers.Print()
	},
}
