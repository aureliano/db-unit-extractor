package reader_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
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
