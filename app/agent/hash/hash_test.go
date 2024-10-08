package hash

import (
	"testing"

	"github.com/AlexBlackNn/metrics/internal/config/configagent"
)

func BenchmarkMetricHash(b *testing.B) {
	cfg := &configagent.Config{
		HashKey: "secret-key",
	}
	hashCalculator := New(cfg)
	inputJSON := `{"id":"%s", "type":"%s", "value": %v}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashCalculator.MetricHash(inputJSON)
	}
}
