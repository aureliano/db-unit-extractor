package reader_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewDataSource(t *testing.T) {
	ds := reader.NewDataSource()
	assert.Equal(t, "", ds.DSN)
	assert.Equal(t, 1, ds.MaxOpenConn)
	assert.Equal(t, 1, ds.MaxIdleConn)
	assert.Nil(t, ds.DB)
}

func TestConnect(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DSN = "sqlite3://file:test.db?cache=shared&mode=memory"

	_, err := ds.Connect(reader.MaxDBTimeout)
	assert.Nil(t, err)

	assert.True(t, ds.IsConnected())

	_, err = ds.Connect(reader.MaxDBTimeout)
	assert.Nil(t, err)
}

func TestConnectOpenError(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DSN = "test://localhost"

	_, err := ds.Connect(reader.MaxDBTimeout)
	assert.Contains(t, err.Error(), "sql: unknown driver \"test\" (forgotten import?)")

	assert.False(t, ds.IsConnected())
}

func TestConnectPingError(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DSN = "sqlite3://file:test.db?cache=shared&mode=memory"

	_, err := ds.Connect(0)
	assert.Contains(t, err.Error(), "context deadline exceeded")

	assert.False(t, ds.IsConnected())
}

func TestDriverName(t *testing.T) {
	ds := reader.NewDataSource()

	assert.Equal(t, "", ds.DriverName())

	ds.DSN = "sqlite3://file:test.db?cache=shared&mode=memory"
	assert.Equal(t, "sqlite3", ds.DriverName())

	ds.DSN = "oracle://usr:pwd@localhost:1521/dbname"
	assert.Equal(t, "oracle", ds.DriverName())
}

func TestConnectionURL(t *testing.T) {
	ds := reader.NewDataSource()

	assert.Equal(t, "", ds.ConnectionURL())

	ds.DSN = "sqlite3://file:test.db?cache=shared&mode=memory"
	assert.Equal(t, "file:test.db?cache=shared&mode=memory", ds.ConnectionURL())

	ds.DSN = "oracle://usr:pwd@localhost:1521/dbname"
	assert.Equal(t, "oracle://usr:pwd@localhost:1521/dbname", ds.ConnectionURL())
}
