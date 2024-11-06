package cmd

import (
	"net"
	"os"

	"github.com/kujourinka/mhws_beta_server/backend"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "mhws_beta_server listen-ip",
	Run: mainRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func mainRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		os.Exit(1)
	}
	if ip := net.ParseIP(args[0]); ip == nil {
		os.Exit(1)
	}
	e := backend.RegisterHandler()

	e.RunTLS(args[0]+":443", "cert/website.crt", "cert/website.key")
}
