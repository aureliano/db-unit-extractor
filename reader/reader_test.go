package reader_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	_, err := reader.NewReader(reader.DataSource{}, nil)
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}

func TestNewOracleReader(t *testing.T) {
	r, err := reader.NewReader(reader.DataSource{DBMSName: "Oracle"}, nil)
	assert.Nil(t, err)
	assert.IsType(t, reader.OracleReader{}, r)
}
