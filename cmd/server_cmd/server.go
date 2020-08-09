package server_cmd

import (
	"github.com/spf13/cobra"
)

var (
	mode     string
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start server",
		Example: "thinkgo server",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	cobra.OnInitialize(initStartCmd)
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

//initStartCmd 命令行初始化
func initStartCmd() {

}

func setup() {
}

func run() error {

	return nil
}
