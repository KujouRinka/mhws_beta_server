package cmd

import (
	"crypto/tls"
	"github.com/kujourinka/mhws_beta_server/backend"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"os"
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

	subCert, err := backend.GenerateDomainCert()
	if err != nil {
		log.Fatalf("Error loading root certificate and key: %v", err)
	}

	server := http.Server{
		Addr:    args[0] + ":443",
		Handler: e,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*subCert},
		},
	}
	_ = server.ListenAndServeTLS("", "")
}
