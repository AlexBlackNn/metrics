package postgres

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
)

// TempMetric is a template to deserialize data from bytes
// (models.MetricGetter and generic types can't be used here).
type TempMetric struct {
	Type  string
	Name  string
	Value any
}

func (m *TempMetric) GetType() string {
	return m.Type
}

func (m *TempMetric) GetName() string {
	return m.Name
}

func (m *TempMetric) GetValue() any {
	return m.Value
}

func (m *TempMetric) GetStringValue() string {
	switch m.GetType() {
	case configserver.MetricTypeCounter:
		if value, ok := m.GetValue().(float64); ok {
			return fmt.Sprintf("%d", int(value))
		}
		if value, ok := m.GetValue().(int64); ok {
			return fmt.Sprintf("%d", value)
		}
	default:
		return fmt.Sprintf("%g", m.GetValue())
	}
	return ""
}
