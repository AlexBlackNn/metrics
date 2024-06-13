package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotValidMetricValue = errors.New("invalid metric value")
var ErrNotValidMetricType = errors.New("invalid metric type")

// Metric works with collected by an agent metrics
type Metric struct {
	Type  string
	Name  string
	Value any
}

// ConvertValueToString converts metric Value to string or returns error
func (m *Metric) ConvertValueToString() (string, error) {
	switch reflect.TypeOf(m.Value).Kind() {
	case reflect.Float64:
		return fmt.Sprintf("%f", m.Value), nil
	case reflect.Uint32:
		return fmt.Sprintf("%d", m.Value), nil
	case reflect.Uint64:
		return fmt.Sprintf("%d", m.Value), nil
	case reflect.Int64:
		return fmt.Sprintf("%d", m.Value), nil
	default:
		return "", fmt.Errorf("unsupported type: %T", m.Value)
	}
}

// Load loads data to metric
func Load(metricType string, metricName string, metricValue string) (Metric, error) {
	var value interface{}
	var err error

	if metricType == "gauge" {
		value, err = strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return Metric{}, ErrNotValidMetricValue
		}
	} else if metricType == "counter" {
		value, err = strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			return Metric{}, ErrNotValidMetricValue
		}
	} else {
		return Metric{}, ErrNotValidMetricType
	}
	return Metric{
		Type:  metricType,
		Name:  strings.ToLower(metricName),
		Value: value,
	}, nil
}
