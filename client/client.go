package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/observiq/bindplane-op-action/client/config"
	"github.com/observiq/bindplane-op-action/client/version"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const KeyHeader = "X-Bindplane-Api-Key"

type BindPlane struct {
	logger *zap.Logger
	config *config.Config
	client *resty.Client
}

// NewBindPlane takes a config and logger and returns a configured BindPlane client
func NewBindPlane(config *config.Config, logger *zap.Logger) (*BindPlane, error) {
	restryClient := resty.New()
	restryClient.SetDisableWarn(true)
	restryClient.SetTimeout(time.Second * 20)

	restryClient.SetBasicAuth(config.Auth.Username, config.Auth.Password)

	if config.Auth.APIKey != "" {
		restryClient.SetHeader(KeyHeader, config.Auth.APIKey)
	}

	restryClient.SetBaseURL(fmt.Sprintf("%s/v1", config.Network.RemoteURL))

	tlsConfig := &tls.Config{}
	if len(config.Network.CertificateAuthority) > 0 {
		tlsConfig.RootCAs = x509.NewCertPool()
		for _, ca := range config.Network.CertificateAuthority {
			if ok := tlsConfig.RootCAs.AppendCertsFromPEM([]byte(ca)); !ok {
				return nil, fmt.Errorf("failed to append certificate authority")
			}
		}
	}

	restryClient.SetTLSClientConfig(tlsConfig)

	return &BindPlane{
		logger: logger,
		config: config,
		client: restryClient,
	}, nil
}

// Version queries the BindPlane API for the version information
func (b *BindPlane) Version(_ context.Context) (version.Version, error) {
	v := version.Version{}
	_, err := b.client.R().SetResult(&v).Get("/version")
	return v, err
}