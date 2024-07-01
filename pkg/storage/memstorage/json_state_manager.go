package memstorage

import (
	"bufio"
	"encoding/json"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"io"
	"os"
)

// dataBaseJSONStateManager saves and restores database state
type dataBaseJSONStateManager struct {
	cfg *configserver.Config
	db  dataBase
}

func (jm *dataBaseJSONStateManager) saveMetrics() error {
	file, err := os.OpenFile(
		jm.cfg.ServerFileStoragePath, os.O_WRONLY|os.O_CREATE, 0777,
	)
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

func (jm *dataBaseJSONStateManager) restoreMetrics() error {
	file, err := os.OpenFile(
		jm.cfg.ServerFileStoragePath, os.O_RDONLY, 0777,
	)
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
