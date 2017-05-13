package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var from string
var to string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tcp-port-forward",
	Short: "a tcp port-forwarder written in golang",
	Long:  `a tcp port-forwarder written in golang`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&from, "from", "f", "localhost:5000", "the local address of the proxy")
	RootCmd.PersistentFlags().StringVarP(&to, "to", "t", "localhost:10000", "the upstream address")
}
