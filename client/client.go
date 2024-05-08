package client

import (
	"context"

	"github.com/observiq/bindplane-op-action/client/config"

	"go.uber.org/zap"
)

type BindPlane struct {
}

func NewBindPlane(config *config.Config, logger *zap.Logger) (*BindPlane, error) {
	return nil, nil
}

func (b *BindPlane) Version(_ context.Context) (string, error) {
	return "", nil
}
