package reader_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDataSource struct{ mock.Mock }

func (ds *mockDataSource) Connect(timeout time.Duration) (*sql.DB, error) {
	args := ds.Called(timeout)
	var db *sql.DB
	if args.Get(0) != nil {
		db, _ = args.Get(0).(*sql.DB)
	}

	return db, args.Error(0)
}

func (ds *mockDataSource) IsConnected() bool {
	args := ds.Called()
	return args.Get(0).(bool)
}

func (ds *mockDataSource) DriverName() string {
	args := ds.Called()
	return args.Get(0).(string)
}

func (ds *mockDataSource) ConnectionURL() string {
	args := ds.Called()
	return args.Get(0).(string)
}

func TestNewReader(t *testing.T) {
	_, err := reader.NewReader(&reader.DataSource{})
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}

func TestNewOracleReader(t *testing.T) {
	ds := new(mockDataSource)
	ds.On("Connect", reader.MaxDBTimeout).Return(nil)
	ds.On("DriverName").Return("oracle")
	ds.On("ConnectionURL").Return("oracle://usr:pwd@localhost:1521/dbname")

	prof := filepath.Join(os.TempDir(), "db.prof")
	defer os.Remove(prof)
	t.Setenv("DB_PROFILE", prof)

	r, err := reader.NewReader(ds)
	assert.Nil(t, err)
	assert.IsType(t, reader.OracleReader{}, r)

	t.Setenv("DB_PROFILE", "")
}

func TestNewOracleReaderError(t *testing.T) {
	ds := new(mockDataSource)
	ds.On("Connect", reader.MaxDBTimeout).Return(fmt.Errorf("connection error"))
	ds.On("DriverName").Return("oracle")
	ds.On("ConnectionURL").Return("oracle://usr:pwd@localhost:1521/dbname")

	_, err := reader.NewReader(ds)
	assert.Equal(t, "connection error", err.Error())
}

func TestNewOracleReaderProfileError(t *testing.T) {
	ds := new(mockDataSource)
	ds.On("Connect", reader.MaxDBTimeout).Return(nil)
	ds.On("DriverName").Return("oracle")
	ds.On("ConnectionURL").Return("oracle://usr:pwd@localhost:1521/dbname")

	t.Setenv("DB_PROFILE", "/path/to/nowhere")

	_, err := reader.NewReader(ds)
	assert.ErrorIs(t, err, os.ErrNotExist)
	t.Setenv("DB_PROFILE", "")
}
