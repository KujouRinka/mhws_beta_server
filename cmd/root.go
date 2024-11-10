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

	var rootCert *x509.Certificate
	var rootPrivateKey *rsa.PrivateKey

	// 检查根证书文件是否存在
	if _, err := os.Stat(rootCertFile); os.IsNotExist(err) {
		// 生成自签名证书
		rootCert, rootPrivateKey, err = generateSelfSignedCert()
		if err != nil {
			log.Fatalf("Error generating self-signed certificate: %v", err)
		}
	} else {
		// 读取根证书和私钥
		rootCert, rootPrivateKey, err = loadRootCertAndKey()
		if err != nil {
			log.Fatalf("Error loading root certificate and key: %v", err)
		}
	}
	subCert, err := generateDomainCert(rootCert, rootPrivateKey)
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

const (
	rootCertFile = "root.crt"
	rootKeyFile  = "root.key"
)

// 生成自签名的根证书
func generateSelfSignedCert() (*x509.Certificate, *rsa.PrivateKey, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// 设置证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"mitmproxy"},
			Country:      []string{"CN"},
			Locality:     []string{"Beijing"},
			Province:     []string{"Beijing"},
			CommonName:   "mitmproxy",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 有效期10年
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 生成自签名证书
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// 将证书编码为PEM格式
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	// 将私钥编码为PEM格式
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// 保存证书到文件
	if err := os.WriteFile(rootCertFile, certPEM, 0644); err != nil {
		return nil, nil, fmt.Errorf("failed to write certificate file: %w", err)
	}

	// 保存私钥到文件
	if err := os.WriteFile(rootKeyFile, privateKeyPEM, 0600); err != nil {
		return nil, nil, fmt.Errorf("failed to write private key file: %w", err)
	}

	return &template, privateKey, nil
}

// 读取根证书和私钥
func loadRootCertAndKey() (*x509.Certificate, *rsa.PrivateKey, error) {
	// 读取证书文件
	certBytes, err := os.ReadFile(rootCertFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// 读取私钥文件
	privateKeyBytes, err := os.ReadFile(rootKeyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// 解析证书
	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("failed to decode PEM certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// 解析私钥
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil || privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return nil, nil, fmt.Errorf("failed to decode PEM private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return cert, privateKey, nil
}

func generateDomainCert(rootCert *x509.Certificate, rootPrivateKey *rsa.PrivateKey) (*tls.Certificate, error) {

	// 生成子证书的私钥
	subjectKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate subject key: %v", err)
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

	// 用根证书和私钥签发子证书
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &template, rootCert, &subjectKey.PublicKey, rootPrivateKey)
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
