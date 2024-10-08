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

func TestMetricHash(t *testing.T) {
	var tests = []struct {
		name       string
		inputData  string
		resultHash string
	}{
		{"simple-json",
			`{"id":"10", "type":"counter", "value": 10}`,
			"ECeoXYtYeQwpMKFMel10cnXJ3MsLTrZFzEANzAAZ/io=",
		},
		{"text",
			`some text to check`,
			"s6oEjfbAbEtNnGwCCl3fFlbHc5sqLsRSmBf3L1wge8o=",
		},
	}
	cfg := &configagent.Config{
		HashKey: "secret-key",
	}
	hashCalculator := New(cfg)
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got := hashCalculator.MetricHash(test.inputData)
			if got != test.resultHash {
				t.Errorf("got %s, want %s", got, test.resultHash)
			}
		})
	}
}
