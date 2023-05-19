package extractor_test

import (
	"os"
	"testing"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DummyConverter string

func (DummyConverter) Convert(_ interface{}, _ *interface{}) {
}

func TestDigestSchemaFileNotFound(t *testing.T) {
	_, err := extractor.DigestSchema("/path/to/unknown/file.yml")
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestDigestSchemaUnmarshalError(t *testing.T) {
	_, err := extractor.DigestSchema("../test/unit/schema_test_unmarshal_error.yml")
	assert.Contains(t, err.Error(), "yaml: unmarshal errors")
}

func TestDigestSchema(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	dataconv.RegisterConverter("conv_timestamp", DummyConverter(""))
	schema, err := extractor.DigestSchema("../test/unit/schema_test.yml")

	require.Nil(t, err)

	converters := schema.Converters
	assert.Len(t, converters, 2)
	assert.EqualValues(t, "conv_date_time", converters[0])
	assert.EqualValues(t, "conv_timestamp", converters[1])

	tables := schema.Tables
	assert.Len(t, tables, 4)

	table := tables[0]
	assert.Equal(t, "customers", table.Name)

	filters := table.Filters
	assert.Len(t, filters, 1)
	assert.EqualValues(t, "id", filters[0].Name)
	assert.EqualValues(t, "1", filters[0].Value)

	columns := table.Columns
	assert.Len(t, columns, 3)
	assert.EqualValues(t, "id", columns[0])
	assert.EqualValues(t, "first_name", columns[1])
	assert.EqualValues(t, "last_name", columns[2])

	ignore := table.Ignore
	assert.Empty(t, ignore)

	table = tables[1]
	assert.Equal(t, "orders", table.Name)

	filters = table.Filters
	assert.Len(t, filters, 1)
	assert.Equal(t, "customer_id", filters[0].Name)
	assert.Equal(t, "1", filters[0].Value)

	columns = table.Columns
	assert.Empty(t, columns)

	ignore = table.Ignore
	assert.Len(t, ignore, 1)
	assert.EqualValues(t, "order_fax", ignore[0])

	table = tables[2]
	assert.Equal(t, "orders_products", table.Name)

	filters = table.Filters
	assert.Len(t, filters, 2)
	assert.EqualValues(t, "order_id", filters[0].Name)
	assert.EqualValues(t, "5", filters[0].Value)
	assert.EqualValues(t, "product_id", filters[1].Name)
	assert.EqualValues(t, "3", filters[1].Value)

	columns = table.Columns
	assert.Empty(t, columns)

	ignore = table.Ignore
	assert.Empty(t, ignore)

	table = tables[3]
	assert.Equal(t, "products", table.Name)

	filters = table.Filters
	assert.Len(t, filters, 1)
	assert.Equal(t, "id", filters[0].Name)
	assert.Equal(t, "3", filters[0].Value)

	columns = table.Columns
	assert.Len(t, columns, 3)
	assert.EqualValues(t, "name", columns[0])
	assert.EqualValues(t, "description", columns[1])
	assert.EqualValues(t, "price", columns[2])

	ignore = table.Ignore
	assert.Empty(t, ignore)
}

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

func TestClassifyGroupOneNotClassified(t *testing.T) {
	s := extractor.Schema{
		Tables: []extractor.TableSchema{
			{Name: "t1", Filters: []extractor.FilterSchema{{Name: "id", Value: "${table.column}"}}},
		},
	}

	err := s.Classify()
	assert.ErrorIs(t, err, extractor.ErrTableClassification)
	assert.Contains(t, err.Error(), "couldn't find any level one tables")
}

func TestClassifyReferenceNotFound(t *testing.T) {
	s := extractor.Schema{
		Tables: []extractor.TableSchema{
			{Name: "t1", Filters: []extractor.FilterSchema{{Name: "id", Value: "1"}}},
			{Name: "t2", Filters: []extractor.FilterSchema{{Name: "id", Value: "${table.column}"}}},
		},
	}

	err := s.Classify()
	assert.ErrorIs(t, err, extractor.ErrTableClassification)
	assert.Contains(t, err.Error(), "t2.id points to unresolvable reference '${table.column}'")
}

func TestClassify(t *testing.T) {
	schema, err := extractor.DigestSchema("../test/unit/schema_test_grouping.yml")
	require.Nil(t, err)

	err = schema.Classify()
	assert.Nil(t, err)

	assert.Equal(t, 1, schema.Tables[1].GroupID)
	assert.Equal(t, 1, schema.Tables[3].GroupID)
	assert.Equal(t, 1, schema.Tables[7].GroupID)

	assert.Equal(t, 2, schema.Tables[0].GroupID)
	assert.Equal(t, 2, schema.Tables[5].GroupID)
	assert.Equal(t, 2, schema.Tables[6].GroupID)
	assert.Equal(t, 2, schema.Tables[8].GroupID)

	assert.Equal(t, 3, schema.Tables[2].GroupID)
	assert.Equal(t, 3, schema.Tables[9].GroupID)
	assert.Equal(t, 3, schema.Tables[12].GroupID)

	assert.Equal(t, 4, schema.Tables[10].GroupID)
	assert.Equal(t, 4, schema.Tables[11].GroupID)

	assert.Equal(t, 5, schema.Tables[4].GroupID)
}
