package cmd

import (
	"github.com/felipejfc/tcp-port-forward/proxy"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "the local port of the port-forwarder",
	Long:  `the local port of the port-forwarder`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy := proxy.NewLocalProxy(from, to)
		err := proxy.Start()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(localCmd)
}
