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
