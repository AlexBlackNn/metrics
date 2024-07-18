package agent

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	restagentsender "github.com/AlexBlackNn/metrics/app/agent/restagentsender/v2"
	"github.com/AlexBlackNn/metrics/internal/config/configagent"
	"log/slog"
)

type CollectSender interface {
	Collect(ctx context.Context)
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
) *AppMonitor {

	hashCalculator := hmac.New(sha256.New, []byte(cfg.HashKey))
	metricsService := restagentsender.New(
		log,
		cfg,
		hashCalculator,
	)
	return &AppMonitor{MetricsService: metricsService}
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
