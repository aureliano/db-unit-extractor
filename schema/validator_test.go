package schema_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
)

func TestValidateSchemaConverter(t *testing.T) {
	s := schema.Model{
		Converters: []schema.Converter{"???"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "converter '???' not found")
}

func TestValidateSchemaNoTableProvided(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := schema.Model{
		Converters: []schema.Converter{"dummy"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "no table provided")
}

func TestValidateSchemaInvalidTable(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := schema.Model{
		Converters: []schema.Converter{"dummy"},
		Tables:     []schema.Table{{Name: "z"}},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'z' invalid name")
}

func TestValidateSchemaRepeatedConverter(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))
	dataconv.RegisterConverter("test1", DummyConverter(""))
	dataconv.RegisterConverter("test2", DummyConverter(""))
	dataconv.RegisterConverter("dummy", DummyConverter(""))
	dataconv.RegisterConverter("test3", DummyConverter(""))

	s := schema.Model{
		Converters: []schema.Converter{"dummy", "test1", "test2", "dummy", "test3"},
		Tables:     []schema.Table{{Name: "tbl"}},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "repeated converter 'dummy'")
}

func TestValidateSchema(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := schema.Model{
		Converters: []schema.Converter{"dummy"},
		Tables:     []schema.Table{{Name: "tbl"}},
	}
	err := s.Validate()
	assert.Nil(t, err)
}

func TestTableSchemaValidate(t *testing.T) {
	s := schema.Table{Name: "tbl_1"}
	err := s.Validate()
	assert.Nil(t, err)
}

func TestTableSchemaValidateName(t *testing.T) {
	s := schema.Table{}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")
}

func TestTableSchemaValidateColumnsAndIgnoreProvided(t *testing.T) {
	s := schema.Table{
		Name:    "tbl",
		Columns: []schema.Column{"a1"},
		Ignore:  []schema.Ignore{"b1"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' with columns and ignore set (excludents)")
}

func TestTableSchemaValidateColumns(t *testing.T) {
	s := schema.Table{
		Name:    "tbl",
		Columns: []schema.Column{"a", "b1"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' validation: 'a' invalid name")
}

func TestTableSchemaValidateIgnore(t *testing.T) {
	s := schema.Table{
		Name:   "tbl",
		Ignore: []schema.Ignore{"a", "b1"},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' validation: 'a' invalid name")
}

func TestTableSchemaValidateFilter(t *testing.T) {
	s := schema.Table{
		Name:    "tbl",
		Filters: []schema.Filter{{Name: "x", Value: "value"}},
	}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "table 'tbl' validation: 'x' invalid name")
}

func TestFilterSchemaValidate(t *testing.T) {
	s := schema.Filter{Name: "", Value: "1"}
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")

	s = schema.Filter{Name: "name", Value: ""}
	err = s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "empty filter value 'name'")

	s = schema.Filter{Name: "name", Value: "value"}
	err = s.Validate()
	assert.Nil(t, err)
}

func TestConverterSchemaValidate(t *testing.T) {
	dataconv.RegisterConverter("dummy", DummyConverter(""))

	s := schema.Converter("")
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "converter '' not found")

	s = schema.Converter("dummy")
	err = s.Validate()
	assert.Nil(t, err)
}

func TestColumnSchemaValidate(t *testing.T) {
	s := schema.Column("")
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")

	s = schema.Column("2fs")
	err = s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'2fs' invalid name")

	s = schema.Column("_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_")
	err = s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "invalid name size")

	s = schema.Column("fs2")
	err = s.Validate()
	assert.Nil(t, err)

	s = schema.Column("_2fs")
	err = s.Validate()
	assert.Nil(t, err)
}

func TestIgnoreSchema(t *testing.T) {
	s := schema.Ignore("")
	err := s.Validate()
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
	assert.Contains(t, err.Error(), "'' invalid name")
}
