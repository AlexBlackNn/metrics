package memstorage

import (
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
)

// TempMetric is a template to deserialize data from bytes
// models.MetricInteraction and generic types can't be used here
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
	case "counter":
		if value, ok := m.GetValue().(float64); ok {
			return fmt.Sprintf("%d", int(value))
		}
		if value, ok := m.GetValue().(int64); ok {
			return fmt.Sprintf("%d", value)
		}
	default:
		return fmt.Sprintf("%f", m.GetValue())
	}
	return ""
}

func (m *TempMetric) AddValue(other models.MetricInteraction) error {
	return errors.New("addValue is not implemented")
}
