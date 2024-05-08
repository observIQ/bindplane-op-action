package client

import (
	"context"

	"github.com/observiq/bindplane-op-action/client/config"
	"github.com/observiq/bindplane-op-action/client/version"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type BindPlane struct {
	logger *zap.Logger
	config *config.Config
	client *resty.Client
}

func NewBindPlane(config *config.Config, logger *zap.Logger) (*BindPlane, error) {
	return &BindPlane{
		logger: logger,
		config: config,
	}, nil
}

// Version queries the BindPlane API for the version information
func (b *BindPlane) Version(_ context.Context) (version.Version, error) {
	v := version.Version{}
	_, err := b.client.R().SetResult(&v).Get("/version")
	return v, err
}
