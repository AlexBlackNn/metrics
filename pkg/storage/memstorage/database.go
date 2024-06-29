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
	fmt.Println(string(jsonData))
	return jsonData, nil
}

func (db *DataBase) decode(data []byte) error {
	fmt.Println("--------", data)
	fmt.Println("--------", string(data))
	return json.Unmarshal(data, &db)
}
