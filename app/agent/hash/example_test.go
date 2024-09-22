package hash

import (
	"fmt"

	"github.com/AlexBlackNn/metrics/internal/config/configagent"
)

// ExampleMetricHash demonstrates the usage of MetricHash.
func ExampleMetricHash() {
	// get config
	cfg, err := configagent.New()
	if err != nil {
		panic(err)
	}
	cfg.HashKey = "secret_key"
	// create metric hasher
	metricHasher := New(cfg)

	// metric_example
	metric := `{
    "delta": 1,
    "id": "string",
    "type": "gauge",
    "value": 11
   }`

	// calculate hash
	hash := metricHasher.MetricHash(metric)
	fmt.Println(hash) // Output: 0r983UzuFKBElN9Pj8e18XUaHzrZRhEGcjLLUPdV/cs=
}
