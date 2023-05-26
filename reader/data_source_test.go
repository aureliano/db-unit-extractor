package reader_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewDataSource(t *testing.T) {
	ds := reader.NewDataSource()
	assert.Equal(t, "", ds.DBMSName)
	assert.Equal(t, "", ds.Username)
	assert.Equal(t, "", ds.Password)
	assert.Equal(t, "", ds.Database)
	assert.Equal(t, "", ds.Hostname)
	assert.Equal(t, 0, ds.Port)
	assert.Equal(t, 1, ds.MaxOpenConn)
	assert.Equal(t, 1, ds.MaxIdleConn)
}

func TestDSNameSqlite3(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DBMSName = "sqlite3"

	expected := "file:test.db?cache=shared&mode=memory"
	actual := ds.DSName()

	assert.Equal(t, expected, actual)
}

func TestDSNameOracle(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DBMSName = "oracle"
	ds.Username = "usr"
	ds.Password = "pwd"
	ds.Hostname = "localhost"
	ds.Port = 1521
	ds.Database = "db_name"

	expected := "oracle://usr:pwd@localhost:1521/db_name"
	actual := ds.DSName()

	assert.Equal(t, expected, actual)
}

func TestConnect(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DBMSName = "sqlite3"

	err := ds.Connect(reader.MaxDBTimeout)
	assert.Nil(t, err)

	assert.True(t, ds.IsConnected())

	err = ds.Connect(reader.MaxDBTimeout)
	assert.Nil(t, err)
}

func TestConnectOpenError(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DBMSName = "test"

	err := ds.Connect(reader.MaxDBTimeout)
	assert.Contains(t, err.Error(), "sql: unknown driver \"test\" (forgotten import?)")

	assert.False(t, ds.IsConnected())
}

func TestConnectPingError(t *testing.T) {
	ds := reader.NewDataSource()
	ds.DBMSName = "sqlite3"

	err := ds.Connect(0)
	assert.Contains(t, err.Error(), "context deadline exceeded")

	assert.False(t, ds.IsConnected())
}
