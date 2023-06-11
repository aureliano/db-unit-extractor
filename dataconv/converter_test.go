package dataconv_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/stretchr/testify/assert"
)

type DummyConverter string

func (DummyConverter) Convert(interface{}) (interface{}, error) {
	return 0, nil
}

func (DummyConverter) Handle(interface{}) bool {
	return true
}

func TestConverterExists(t *testing.T) {
	assert.False(t, dataconv.ConverterExists(""))
}

func TestRegisterConverter(t *testing.T) {
	dataconv.RegisterConverter("xpto", DummyConverter(""))
	assert.True(t, dataconv.ConverterExists("xpto"))
}

func TestRegisterConverters(t *testing.T) {
	dataconv.RegisterConverters()
	assert.True(t, dataconv.ConverterExists(dataconv.DateTimeISO8601ID))
	assert.True(t, dataconv.ConverterExists(dataconv.BlobConverterID))
}

func TestGetConverter(t *testing.T) {
	dataconv.RegisterConverters()

	c := dataconv.GetConverter(dataconv.DateTimeISO8601ID)
	assert.NotNil(t, c)

	c = dataconv.GetConverter(dataconv.BlobConverterID)
	assert.NotNil(t, c)
}
