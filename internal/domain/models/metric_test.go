package models

import (
	"testing"

	"github.com/AlexBlackNn/metrics/internal/config/configserver"
)

func TestNewMetric(t *testing.T) {
	t.Run("Gauge", func(t *testing.T) {
		metric, err := New(configserver.MetricTypeGauge, "test_gauge", "12.34")
		if err != nil {
			t.Errorf("Error creating gauge metric: %v", err)
		}

		gaugeMetric, ok := metric.(*Metric[float64])
		if !ok {
			t.Errorf("Expected metric to be of type *Metric[float64], got %T", metric)
		}

		if gaugeMetric.Type != configserver.MetricTypeGauge {
			t.Errorf("Expected metric type to be %s, got %s", configserver.MetricTypeGauge, gaugeMetric.Type)
		}

		if gaugeMetric.Name != "test_gauge" {
			t.Errorf("Expected metric name to be test_gauge, got %s", gaugeMetric.Name)
		}

		if gaugeMetric.Value != 12.34 {
			t.Errorf("Expected metric value to be 12.34, got %v", gaugeMetric.Value)
		}
	})

	t.Run("Counter", func(t *testing.T) {
		metric, err := New(configserver.MetricTypeCounter, "test_counter", "1234")
		if err != nil {
			t.Errorf("Error creating counter metric: %v", err)
		}

		counterMetric, ok := metric.(*Metric[uint64])
		if !ok {
			t.Errorf("Expected metric to be of type *Metric[uint64], got %T", metric)
		}

		if counterMetric.Type != configserver.MetricTypeCounter {
			t.Errorf("Expected metric type to be %s, got %s", configserver.MetricTypeCounter, counterMetric.Type)
		}

		if counterMetric.Name != "test_counter" {
			t.Errorf("Expected metric name to be test_counter, got %s", counterMetric.Name)
		}

		if counterMetric.Value != 1234 {
			t.Errorf("Expected metric value to be 1234, got %v", counterMetric.Value)
		}
	})

	t.Run("InvalidType", func(t *testing.T) {
		_, err := New("invalid_type", "test_metric", "1234")
		if err == nil {
			t.Errorf("Expected error creating metric with invalid type, got nil")
		}
	})

	t.Run("InvalidValue", func(t *testing.T) {
		_, err := New(configserver.MetricTypeGauge, "test_gauge", "invalid_value")
		if err == nil {
			t.Errorf("Expected error creating metric with invalid value, got nil")
		}
	})
}

func TestMetric_AddValue(t *testing.T) {

	t.Run("AddCounter", func(t *testing.T) {
		metric1, _ := New(configserver.MetricTypeCounter, "test_counter", "1234")
		metric2, _ := New(configserver.MetricTypeCounter, "test_counter", "567")

		err := metric1.AddValue(metric2)
		if err != nil {
			t.Errorf("Error adding counter metrics: %v", err)
		}

		counterMetric, ok := metric1.(*Metric[uint64])
		if !ok {
			t.Errorf("Expected metric1 to be of type *Metric[uint64], got %T", metric1)
		}

		if counterMetric.Value != 1801 {
			t.Errorf("Expected metric value to be 1801, got %v", counterMetric.Value)
		}
	})

	t.Run("DifferentTypes", func(t *testing.T) {
		metric1, _ := New(configserver.MetricTypeGauge, "test_gauge", "12.34")
		metric2, _ := New(configserver.MetricTypeCounter, "test_gauge", "567")

		err := metric1.AddValue(metric2)
		if err == nil {
			t.Errorf("Expected error adding metrics of different types, got nil")
		}
	})

	t.Run("DifferentNames", func(t *testing.T) {
		metric1, _ := New(configserver.MetricTypeGauge, "test_gauge1", "12.34")
		metric2, _ := New(configserver.MetricTypeGauge, "test_gauge2", "5.67")

		err := metric1.AddValue(metric2)
		if err == nil {
			t.Errorf("Expected error adding metrics of different names, got nil")
		}
	})
}

func TestMetric_GetStringValue(t *testing.T) {
	t.Run("Gauge", func(t *testing.T) {
		metric, _ := New(configserver.MetricTypeGauge, "test_gauge", "12.34")
		stringValue := metric.GetStringValue()
		if stringValue != "12.340000" {
			t.Errorf("Expected string value to be 12.340000, got %s", stringValue)
		}
	})

	t.Run("Counter", func(t *testing.T) {
		metric, _ := New(configserver.MetricTypeCounter, "test_counter", "1234")
		stringValue := metric.GetStringValue()
		if stringValue != "1234" {
			t.Errorf("Expected string value to be 1234, got %s", stringValue)
		}
	})
}

func TestCheckModelType(t *testing.T) {
	t.Run("ValidType", func(t *testing.T) {
		err := CheckModelType(configserver.MetricTypeGauge)
		if err != nil {
			t.Errorf("Error checking valid metric type: %v", err)
		}
	})

	t.Run("InvalidType", func(t *testing.T) {
		err := CheckModelType("invalid_type")
		if err == nil {
			t.Errorf("Expected error checking invalid metric type, got nil")
		}
	})
}
