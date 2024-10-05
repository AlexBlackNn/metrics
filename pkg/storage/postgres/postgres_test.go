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
