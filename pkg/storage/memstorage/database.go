package memstorage

import (
	"encoding/json"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
)

type MetricData struct {
	Type  string `json:"Type"`
	Name  string `json:"Name"`
	Value any    `json:"Value"`
}

func (m *MetricData) GetType() string {
	return m.Type
}

func (m *MetricData) GetName() string {
	return m.Name
}

func (m *MetricData) GetValue() any {
	return m.Value
}

func (m *MetricData) GetStringValue() string {

	switch m.GetValue().(type) {
	case uint64, uint32:
		return fmt.Sprintf("%d", m.GetValue())
	default:
		return fmt.Sprintf("%f", m.GetValue())
	}
}

// AddValue adds the value of another Metric to the current Metric
func (m *MetricData) AddValue(other models.MetricInteraction) error {
	return nil
}

type DataBase map[string]models.MetricInteraction

func (db *DataBase) encode() ([]byte, error) {
	jsonData, err := json.Marshal(*db)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (db *DataBase) decode(data []byte) error {
	tempDB := make(map[string]MetricData)
	err := json.Unmarshal(data, &tempDB)
	if err != nil {
		return err
	}
	for k, v := range tempDB {
		v := v
		(*db)[k] = &v
	}
	return nil
}
