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

func TestConverterExists(t *testing.T) {
	assert.False(t, dataconv.ConverterExists(""))
}

func TestRegisterConverter(t *testing.T) {
	dataconv.RegisterConverter("xpto", DummyConverter(""))
	assert.True(t, dataconv.ConverterExists("xpto"))
}
