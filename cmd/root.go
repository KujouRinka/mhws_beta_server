package cmd

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/kujourinka/mhws_beta_server/backend"
	"github.com/kujourinka/mhws_beta_server/config"
	"github.com/spf13/cobra"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use: "mhws_beta_server [-a address] [-p port] [-c root_cert] [-k root_key] [--cert-domain \"d1,d2\"] [--api-host host]",
	Run: mainRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	initFlags()
	cobra.OnInitialize(initOthers)
}

func initFlags() {
	rootCmd.PersistentFlags().StringVarP(&cfg.ListenAddr, "address", "a", "0.0.0.0", "specify server host")
	rootCmd.PersistentFlags().Uint16VarP(&cfg.ListenPort, "port", "p", 443, "specify server port")
	rootCmd.PersistentFlags().StringVarP(&cfg.RootCertFile, "root-cert", "c", "./cert/root.crt", "root cert file path")
	rootCmd.PersistentFlags().StringVarP(&cfg.RootKeyFile, "root-key", "k", "./cert/root.key", "root key file path")
	rootCmd.PersistentFlags().StringSliceVar(&cfg.CertDomain, "cert-domain", []string{}, "domain name for which certificate is required")
	// rootCmd.PersistentFlags().StringSliceVar(&cfg.CertDomain, "cert-domain", []string{"hjm.rebe.capcom.com", "40912.playfabapi.com"}, "domain name for which certificate is required")
	rootCmd.PersistentFlags().StringVar(&cfg.ApiHost, "api-host", "mhws.io", "api host")
}

func initOthers() {
	cfg.CertDomain = append(cfg.CertDomain, cfg.ApiHost)
	m := map[string]struct{}{}
	ns := make([]string, 0)
	for _, v := range cfg.CertDomain {
		if v == "" {
			continue
		}
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			ns = append(ns, v)
		}
	}
	cfg.CertDomain = ns
}

func mainRun(cmd *cobra.Command, args []string) {
	subCert, err := backend.GenerateDomainCert(&cfg)
	if err != nil {
		log.Fatalf("Error loading root certificate and key: %v", err)
	}

	e := backend.RegisterHandler(&cfg)
	server := http.Server{
		Addr:    cfg.ListenAddr + ":" + strconv.Itoa(int(cfg.ListenPort)),
		Handler: e,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*subCert},
		},
	}
	_ = server.ListenAndServeTLS("", "")
}
