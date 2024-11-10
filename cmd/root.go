package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

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

	subCert, err := createSelfSignedCert("cert/root.crt", "cert/root.key")
	if err != nil {
		log.Fatalf("createSelfSingedCert error: %v", err)
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

func createSelfSignedCert(certFile, keyFile string) (*tls.Certificate, error) {
	// 加载根证书和私钥
	rootCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load root certificate and key: %v", err)
	}

	// 解析根证书
	rootCertParsed, err := x509.ParseCertificate(rootCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse root certificate: %v", err)
	}

	// 解析根私钥
	rootKey, ok := rootCert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("root private key is not an RSA key")
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"mimtproxy"},
			OrganizationalUnit: []string{"IT Department"},
			CommonName:         "example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false, // 子证书不是CA证书
		DNSNames:              []string{"hjm.rebe.capcom.com", "mhws.io", "40912.playfabapi.com"},
	}

	// 生成子证书的私钥
	subjectKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate subject key: %v", err)
	}

	// 用根证书和私钥签发子证书
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &template, rootCertParsed, &subjectKey.PublicKey, rootKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	// 将证书编码为PEM格式
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDERBytes})

	// 将私钥编码为PEM格式
	subjectKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(subjectKey)})

	// 将PEM格式的证书和私钥加载到tls.Certificate结构体中
	subCert, err := tls.X509KeyPair(certPEM, subjectKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load sub certificate and key: %v", err)
	}
	return &subCert, nil
}
