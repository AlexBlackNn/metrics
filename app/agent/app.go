package agent

import (
	"context"
	"log/slog"

	"github.com/AlexBlackNn/metrics/app/agent/encryption"
	restagentsender "github.com/AlexBlackNn/metrics/app/agent/restagentsender/v2"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
)

type CollectSender interface {
	Collect(ctx context.Context)
	CollectAddition(ctx context.Context)
	Send(ctx context.Context)
}

// AppMonitor service consists all service layers.
type AppMonitor struct {
	MetricsService CollectSender
}

// NewAppMonitor creates App.
func NewAppMonitor(
	log *slog.Logger,
	cfg *configagent.Config,
) (*AppMonitor, error) {

	encryptor, err := encryption.NewEncryptor(cfg.CryptoKeyPath)
	if err != nil {
		return nil, err
	}

	metricsService := restagentsender.New(
		log,
		cfg,
		encryptor,
	)
	return &AppMonitor{MetricsService: metricsService}, nil
}
