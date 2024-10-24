package postgres

import (
	"testing"
)

func TestMetricHash(t *testing.T) {
	var tests = []struct {
		metricType        string
		metricName        string
		metricValue       any
		metricStringValue string
	}{
		{"counter",
			"test_counter",
			int64(10),
			"10",
		},
		{"counter",
			"test_counter",
			float64(10),
			"10",
		},
		{"gauge",
			"test_gauge",
			0.21,
			"0.21",
		},
	}

	for _, test := range tests {
		t.Run(test.metricName, func(t *testing.T) {
			tempMetric := TempMetric{
				test.metricType,
				test.metricName,
				test.metricValue,
			}
			metricTypeGot := tempMetric.GetType()
			if metricTypeGot != test.metricType {
				t.Errorf("got %s, want %s", metricTypeGot, test.metricType)
			}
			metricNameGot := tempMetric.GetName()
			if metricNameGot != test.metricName {
				t.Errorf("got %s, want %s", metricNameGot, test.metricName)
			}
			metricValueGot := tempMetric.GetValue()
			if metricValueGot != test.metricValue {
				t.Errorf("got %v, want %v", metricValueGot, test.metricType)
			}
			metricStringValueGot := tempMetric.GetStringValue()
			if metricStringValueGot != test.metricStringValue {
				t.Errorf("got %v, want %v", metricStringValueGot, test.metricStringValue)
			}
		})
	}
}
