package config

type Config struct {
	ListenAddr string
	ListenPort uint16

	RootCertFile string
	RootKeyFile  string
	CertDomain   []string

	ApiHost string
}
