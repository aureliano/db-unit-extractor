package reader_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	_, err := reader.NewReader(&reader.DataSource{})
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}

func TestNewOracleReader(t *testing.T) {
	r, err := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname"})
	assert.Nil(t, err)
	assert.IsType(t, reader.OracleReader{}, r)
}
