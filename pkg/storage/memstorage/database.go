package memstorage

import (
	"encoding/json"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
)

type DataBase map[string]models.MetricInteraction

func (db *DataBase) encode() ([]byte, error) {
	jsonData, err := json.Marshal(*db)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (db *DataBase) UnmarshalJSON(data []byte) error {

	var TempDBMetric map[string]TempMetric

	if err := json.Unmarshal(data, &TempDBMetric); err != nil {
		return err
	}
	for _, v := range TempDBMetric {
		v := v
		fmt.Println(v.Type, v.Name, v.GetStringValue())
		metric, err := models.New(v.Type, v.Name, v.GetStringValue())
		if err != nil {
			return err
		}
		(*db)[v.Name] = metric
	}
	return nil
}

func (db *DataBase) decode(data []byte) error {
	err := json.Unmarshal(data, &db)
	if err != nil {
		return err
	}
	return nil
}
