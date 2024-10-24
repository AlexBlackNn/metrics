package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

var testDbInstance *sql.DB

func TestPostStorage(t *testing.T) {
	ds := &PostStorage{NewTemplate(), testDbInstance}

	tests := []struct {
		name     string
		testFunc func(*testing.T, *PostStorage)
	}{
		{"CreateConnection", testCreateConnection},
		{"HealthCheck", testHealthCheck},
		{"UpdateMetric", testUpdateMetric},
		{"UpdateSeveralMetrics", testUpdateSeveralMetrics},
		{"GetAllMetrics", testGetAllMetrics},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, ds)
		})
	}
}

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func testCreateConnection(t *testing.T, ds *PostStorage) {
	assert.NotNil(t, ds)
}

func testHealthCheck(t *testing.T, ds *PostStorage) {
	err := ds.HealthCheck(context.Background())
	assert.NoError(t, err)
}

func testUpdateMetric(t *testing.T, ds *PostStorage) {
	testMetric := &models.Metric[uint64]{
		Type:  "counter",
		Name:  "test_metric_counter",
		Value: 10,
	}

	err := ds.UpdateMetric(context.Background(), testMetric)

	assert.NoError(t, err)

	testMetricGot, err := ds.GetMetric(context.Background(), testMetric)
	assert.Equal(t, testMetric.GetStringValue(), testMetricGot.GetStringValue())
	assert.Equal(t, testMetric.GetName(), testMetricGot.GetName())
	assert.Equal(t, testMetric.GetValue(), testMetricGot.GetValue())
	assert.NoError(t, err)
}

func testUpdateSeveralMetrics(t *testing.T, ds *PostStorage) {
	metrics := make(map[string]models.MetricGetter)

	metrics["testMetric1"] = &models.Metric[uint64]{
		Type:  "counter",
		Name:  "testMetric1",
		Value: 10,
	}

	metrics["testMetric2"] = &models.Metric[float64]{
		Type:  "gauge",
		Name:  "testMetric2",
		Value: 0.10,
	}

	err := ds.UpdateSeveralMetrics(context.Background(), metrics)

	assert.NoError(t, err)

	testMetricGot, err := ds.GetMetric(context.Background(), metrics["testMetric1"])
	assert.NoError(t, err)
	assert.Equal(t, metrics["testMetric1"].GetStringValue(), testMetricGot.GetStringValue())
	assert.Equal(t, metrics["testMetric1"].GetName(), testMetricGot.GetName())
	assert.Equal(t, metrics["testMetric1"].GetValue(), testMetricGot.GetValue())

	testMetricGot, err = ds.GetMetric(context.Background(), metrics["testMetric2"])
	assert.NoError(t, err)
	assert.Equal(t, metrics["testMetric2"].GetStringValue(), testMetricGot.GetStringValue())
	assert.Equal(t, metrics["testMetric2"].GetName(), testMetricGot.GetName())
	assert.Equal(t, metrics["testMetric2"].GetValue(), testMetricGot.GetValue())
}

func testGetAllMetrics(t *testing.T, ds *PostStorage) {
	metrics := make(map[string]models.MetricGetter)

	metrics["testMetric1"] = &models.Metric[uint64]{
		Type:  "counter",
		Name:  "testMetric1",
		Value: 10,
	}

	metrics["testMetric2"] = &models.Metric[float64]{
		Type:  "gauge",
		Name:  "testMetric2",
		Value: 0.10,
	}

	err := ds.UpdateSeveralMetrics(context.Background(), metrics)
	assert.NoError(t, err)
	metricsGot, err := ds.GetAllMetrics(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, metricsGot)
}
