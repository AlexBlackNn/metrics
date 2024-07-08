package memstorage

import (
	"encoding/json"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
)

type dataBase map[string]models.MetricGetter

func (db *dataBase) UnmarshalJSON(data []byte) error {

	// Can't unmarshal to models.MetricInteraction (interface).
	var TempDBMetric map[string]TempMetric

	if err := json.Unmarshal(data, &TempDBMetric); err != nil {
		return err
	}
	for _, v := range TempDBMetric {
		v := v
		metric, err := models.New(v.Type, v.Name, v.GetStringValue())
		if err != nil {
			return err
		}
		(*db)[v.Name] = metric
	}
	return nil
}
