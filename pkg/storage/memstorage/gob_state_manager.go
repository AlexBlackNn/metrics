package memstorage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"os"
)

func init() {
	gob.Register(models.Metric[uint64]{})
	gob.Register(models.Metric[float64]{})
	gob.Register(encodeMetricUint64)
	gob.Register(encodeMetricFloat64)

}

// dataBaseJSONStateManager saves and restores database state
type dataBaseGOBStateManager struct {
	cfg *configserver.Config
	db  dataBase
}

// Custom gob encoder for models.Metric[uint64]
func encodeMetricUint64(enc *gob.Encoder, m models.Metric[uint64]) error {
	return enc.Encode(struct {
		Name  string
		Type  string
		Value uint64
	}{
		Name:  m.Name,
		Type:  m.Type,
		Value: m.Value,
	})
}

// Custom gob encoder for models.Metric[uint64]
func encodeMetricFloat64(enc *gob.Encoder, m models.Metric[float64]) error {
	return enc.Encode(struct {
		Name  string
		Type  string
		Value float64
	}{
		Name:  m.Name,
		Type:  m.Type,
		Value: m.Value,
	})
}

func (gm *dataBaseGOBStateManager) saveMetrics() error {

	file, err := os.OpenFile(
		gm.cfg.ServerFileStoragePath, os.O_WRONLY|os.O_CREATE, 0777,
	)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	var buffer bytes.Buffer
	// кодирование
	if err = gob.NewEncoder(&buffer).Encode(gm.db); err != nil {
		fmt.Println(err)
		return err
	}
	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (gm *dataBaseGOBStateManager) restoreMetrics() error {
	file, err := os.OpenFile(
		gm.cfg.ServerFileStoragePath, os.O_RDONLY, 0777,
	)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := bufio.NewReader(file)
	dec := gob.NewDecoder(reader) // Create a gob decoder

	// Decode the data into a map[string]interface{}
	var decodedData map[string]interface{}
	if err = dec.Decode(&decodedData); err != nil {
		return err
	}
	for k, v := range decodedData {
		switch v := v.(type) {
		case models.Metric[uint64]:
			gm.db[k] = &v
		case models.Metric[float64]:
			gm.db[k] = &v
		default:
			return errors.New("unknown type")
		}
	}
	return nil
}