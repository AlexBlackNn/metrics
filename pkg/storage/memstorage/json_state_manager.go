package memstorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"io"
	"log/slog"
	"os"
	"sync"
)

// dataBaseJSONStateManager saves and restores database state.
type dataBaseJSONStateManager struct {
	cfg   *configserver.Config
	log   *slog.Logger
	db    dataBase
	mutex *sync.RWMutex
}

func (jm *dataBaseJSONStateManager) saveMetrics() error {
	log := jm.log.With(
		slog.String("info", "STORAGE LAYER: json_state_manager.saveMetrics"),
	)
	log.Debug("starts saving metric")
	jm.mutex.RLock()
	defer jm.mutex.RUnlock()
	file, err := os.OpenFile(
		jm.cfg.ServerFileStoragePath, os.O_WRONLY|os.O_CREATE, 0777,
	)
	if err != nil {
		return fmt.Errorf(
			"STORAGE LAYER: json_state_manager.saveMetrics: couldn't open metric file: %w",
			err,
		)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error("failed to close file", "err", err)
		}
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
	log := jm.log.With(
		slog.String("info", "STORAGE LAYER: json_state_manager.restoreMetrics"),
	)
	log.Info("restore saved metric")

	jm.mutex.RLock()
	defer jm.mutex.RUnlock()
	file, err := os.OpenFile(
		jm.cfg.ServerFileStoragePath, os.O_RDONLY, 0777,
	)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error("failed to close file", "err", err)
		}
	}(file)

	reader := bufio.NewReader(file)
	tmpBuffer, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf(
			"STORAGE LAYER: json_state_manager.restoreMetrics: couldn't read metric file: %w",
			err,
		)
	}

	return json.Unmarshal(tmpBuffer, &jm.db)
}
