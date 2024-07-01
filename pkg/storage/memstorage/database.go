package memstorage

import (
	"bufio"
	"encoding/json"
	"github.com/AlexBlackNn/metrics/internal/config"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"io"
	"os"
)

type DataBase map[string]models.MetricInteraction

func (db *DataBase) UnmarshalJSON(data []byte) error {

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

type DataBaseJsonStateManager struct {
	cfg *config.Config
	db  DataBase
}

func (jm *DataBaseJsonStateManager) SaveMetrics() error {
	file, err := os.OpenFile(jm.cfg.ServerFileStoragePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	dataBaseBytes, err := json.Marshal(jm.db)
	if err != nil {
		return err
	}
	_, err = writer.Write(dataBaseBytes)
	if err != nil {
		return err
	}
	return nil
}

func (jm *DataBaseJsonStateManager) RestoreMetrics() error {
	file, err := os.OpenFile(jm.cfg.ServerFileStoragePath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := bufio.NewReader(file)
	tmpBuffer, err := io.ReadAll(reader)
	if err != nil {
		return ErrFailedToRestoreMetrics
	}
	err = json.Unmarshal(tmpBuffer, &jm.db)
	if err != nil {
		return ErrFailedToRestoreMetrics
	}
	return nil
}
