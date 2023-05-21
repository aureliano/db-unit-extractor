package reader_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	_, err := reader.NewReader(reader.DataSource{})
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}