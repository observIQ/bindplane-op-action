package config

type Config struct {
	Auth    Auth
	Network Network
}

type Auth struct {
	APIKey   string
	Username string
	Password string
}

type Network struct {
	RemoteURL            string
	CertificateAuthority []string
	TLS
}

type TLS struct {
	CertificateAuthority []string
}
