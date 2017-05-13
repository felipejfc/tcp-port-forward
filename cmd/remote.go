package cmd

import (
	"github.com/felipejfc/tcp-port-forward/proxy"
	"github.com/spf13/cobra"
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "init remote part of the port-forwarder",
	Long:  `init the remote part of the port-forwarder`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy := proxy.NewRemoteProxy(from, to)
		err := proxy.Start()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(remoteCmd)
}
