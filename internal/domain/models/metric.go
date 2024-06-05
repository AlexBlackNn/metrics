package models

import (
	"fmt"
	"reflect"
)

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
