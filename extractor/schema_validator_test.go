package extractor_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/stretchr/testify/assert"
)

func TestValidateSchemaConverter(t *testing.T) {
	s := extractor.Schema{
		Converters: []extractor.ConverterSchema{"???"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "converter '???' not found")
}

func TestValidateSchemaNoTableProvided(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := extractor.Schema{
		Converters: []extractor.ConverterSchema{"dummy"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "no table provided")
}

func TestValidateSchemaInvalidTable(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := extractor.Schema{
		Converters: []extractor.ConverterSchema{"dummy"},
		Tables:     []extractor.TableSchema{{Name: "z"}},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'z' invalid name")
}

func TestValidateSchema(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := extractor.Schema{
		Converters: []extractor.ConverterSchema{"dummy"},
		Tables:     []extractor.TableSchema{{Name: "tbl"}},
	}
	err := s.Validate()
	assert.Nil(t, err)
}

func TestTableSchemaValidate(t *testing.T) {
	s := extractor.TableSchema{Name: "tbl_1"}
	err := s.Validate()
	assert.Nil(t, err)
}

func TestTableSchemaValidateName(t *testing.T) {
	s := extractor.TableSchema{}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")
}

func TestTableSchemaValidateColumnsAndIgnoreProvided(t *testing.T) {
	s := extractor.TableSchema{
		Name:    "tbl",
		Columns: []extractor.ColumnSchema{"a1"},
		Ignore:  []extractor.IgnoreSchema{"b1"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' with columns and ignore set (excludents)")
}

func TestTableSchemaValidateColumns(t *testing.T) {
	s := extractor.TableSchema{
		Name:    "tbl",
		Columns: []extractor.ColumnSchema{"a", "b1"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' validation: 'a' invalid name")
}

func TestTableSchemaValidateIgnore(t *testing.T) {
	s := extractor.TableSchema{
		Name:   "tbl",
		Ignore: []extractor.IgnoreSchema{"a", "b1"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' validation: 'a' invalid name")
}

func TestTableSchemaValidateFilter(t *testing.T) {
	s := extractor.TableSchema{
		Name:    "tbl",
		Filters: []extractor.FilterSchema{{Name: "x", Value: "value"}},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' validation: 'x' invalid name")
}

func TestFilterSchemaValidate(t *testing.T) {
	s := extractor.FilterSchema{Name: "", Value: "1"}
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")

	s = extractor.FilterSchema{Name: "name", Value: ""}
	err = s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "empty filter value 'name'")

	s = extractor.FilterSchema{Name: "name", Value: "value"}
	err = s.Validate()
	assert.Nil(t, err)
}

func TestConverterSchemaValidate(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := extractor.ConverterSchema("")
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "converter '' not found")

	s = extractor.ConverterSchema("dummy")
	err = s.Validate()
	assert.Nil(t, err)
}

func TestColumnSchemaValidate(t *testing.T) {
	s := extractor.ColumnSchema("")
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")

	s = extractor.ColumnSchema("2fs")
	err = s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'2fs' invalid name")

	s = extractor.ColumnSchema("_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_")
	err = s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "invalid name size")

	s = extractor.ColumnSchema("fs2")
	err = s.Validate()
	assert.Nil(t, err)

	s = extractor.ColumnSchema("_2fs")
	err = s.Validate()
	assert.Nil(t, err)
}

func TestIgnoreSchema(t *testing.T) {
	s := extractor.IgnoreSchema("")
	err := s.Validate()
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")
}
