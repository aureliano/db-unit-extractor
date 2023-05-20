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

func TestDigestSchemaUnmarshalErrorUnknownField(t *testing.T) {
	_, err := extractor.DigestSchema("../test/unit/schema_test_unmarshal_error_unknown_field.yml")
	assert.Contains(t, err.Error(), "line 7: field ArbitraryField not found in type extractor.TableSchema")
}

func TestDigestSchemaValidationError(t *testing.T) {
	_, err := extractor.DigestSchema("../test/unit/schema_test.yml")
	assert.ErrorIs(t, err, extractor.ErrSchemaValidation)
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
