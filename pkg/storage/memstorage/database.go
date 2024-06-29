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

func (m *MetricData) GetStringValue() string {

	switch m.Type {
	case "counter":
		return fmt.Sprintf("%d", m.Value)
	default:
		return fmt.Sprintf("%f", m.Value)
	}
}

type DataBase map[string]models.MetricInteraction

func (db *DataBase) encode() ([]byte, error) {
	jsonData, err := json.Marshal(*db)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (db DataBase) MarshalJSON() ([]byte, error) {
	// чтобы избежать рекурсии при json.Unmarshal, объявляем новый тип
	var dataMetric []models.MetricInteraction
	for _, v := range db {
		dataMetric = append(dataMetric, v)
	}
	return json.Marshal(dataMetric)
}

func (db *DataBase) UnmarshalJSON(data []byte) error {
	var dataMetric []MetricData

	if err := json.Unmarshal(data, &dataMetric); err != nil {
		fmt.Println("***********", err)
		return err
	}
	for _, v := range dataMetric {
		v := v
		fmt.Println("---------", v.Type, v.Name, v.GetStringValue())
		metric, err := models.New(v.Type, v.Name, v.GetStringValue())
		if err != nil {
			fmt.Println("(((((((((((((", err)
			return err
		}
		fmt.Println(")))))))))))", metric)
		(*db)[v.Name] = metric
	}
	return nil
}

func (db *DataBase) decode(data []byte) error {
	var realDB DataBase
	err := json.Unmarshal(data, &realDB)
	if err != nil {
		return err
	}
	return nil
}
