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

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestCreateConnection(t *testing.T) {
	ds := &PostStorage{
		NewTemplate(),
		testDbInstance,
	}
	assert.NotNil(t, ds)
}

func TestHealthCheck(t *testing.T) {
	ds := &PostStorage{
		NewTemplate(),
		testDbInstance,
	}
	err := ds.HealthCheck(context.Background())
	assert.NoError(t, err)
}

func TestUpdateSeveralMetrics(t *testing.T) {
	ds := &PostStorage{
		NewTemplate(),
		testDbInstance,
	}

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
	assert.Equal(t, metrics["testMetric1"].GetStringValue(), testMetricGot.GetStringValue())
	assert.Equal(t, metrics["testMetric1"].GetName(), testMetricGot.GetName())
	assert.Equal(t, metrics["testMetric1"].GetValue(), testMetricGot.GetValue())
	assert.NoError(t, err)

	testMetricGot, err = ds.GetMetric(context.Background(), metrics["testMetric2"])
	assert.Equal(t, metrics["testMetric2"].GetStringValue(), testMetricGot.GetStringValue())
	assert.Equal(t, metrics["testMetric2"].GetName(), testMetricGot.GetName())
	assert.Equal(t, metrics["testMetric2"].GetValue(), testMetricGot.GetValue())
	assert.NoError(t, err)
}

func TestUpdateMetric(t *testing.T) {
	ds := &PostStorage{
		NewTemplate(),
		testDbInstance,
	}

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
