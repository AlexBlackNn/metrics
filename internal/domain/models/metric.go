package models

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"golang.org/x/exp/constraints"
	"strconv"
)

type MetricAdder interface {
	AddValue(metric MetricGetter) error
}

type MetricGetter interface {
	GetType() string
	GetName() string
	GetValue() any
	GetStringValue() string
}

type MetricInteraction interface {
	MetricAdder
	MetricGetter
}

// Metric works with metrics collected by an agent.
type Metric[T constraints.Integer | constraints.Float] struct {
	Type  string
	Name  string
	Value T
}

func (m *Metric[T]) GetType() string {
	return m.Type
}

func (m *Metric[T]) GetName() string {
	return m.Name
}

func (m *Metric[T]) GetValue() any {
	return m.Value
}

func (m *Metric[T]) GetStringValue() string {

	switch m.GetValue().(type) {
	case uint64, uint32:
		return fmt.Sprintf("%d", m.GetValue())
	default:
		return fmt.Sprintf("%f", m.GetValue())
	}
}

// AddValue adds the value of another Metric to the current Metric.
func (m *Metric[T]) AddValue(other MetricGetter) error {
	if m.GetType() != other.GetType() {
		return ErrAddDifferentMetricType
	}
	if m.GetName() != other.GetName() {
		return ErrAddDifferentMetricName
	}

	// Since T is constrained to be either constraints.Float or constraints.Integer, we can use them here.
	if mValue, ok := any(m.Value).(float64); ok {
		if oValue, ok := other.GetValue().(float64); ok {
			m.Value = T(mValue + oValue)
		}
	}
	if mValue, ok := any(m.Value).(uint64); ok {
		if oValue, ok := other.GetValue().(uint64); ok {
			m.Value = T(mValue + oValue)
			return nil
		}
	}
	return ErrAddMetricValueCast
}

func CheckModelType(metricType string) error {
	if metricType != configserver.MetricTypeGauge && metricType != configserver.MetricTypeCounter {
		return ErrNotValidMetricType
	}
	return nil
}

func New(metricType string, metricName string, metricValue string) (MetricInteraction, error) {

	if metricType == configserver.MetricTypeGauge {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return nil, ErrNotValidMetricValue
		}
		return &Metric[float64]{
			Type:  metricType,
			Name:  metricName,
			Value: value,
		}, nil
	}

	if metricType == configserver.MetricTypeCounter {
		value, err := strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			return nil, ErrNotValidMetricValue
		}
		return &Metric[uint64]{
			Type:  metricType,
			Name:  metricName,
			Value: value,
		}, nil
	}
	return nil, ErrNotValidMetricType
}
