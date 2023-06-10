package schema_test

import (
	"os"
	"testing"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DummyConverter string

func (DummyConverter) Convert(interface{}) (interface{}, error) {
	return 0, nil
}

func (DummyConverter) Handle(interface{}) bool { return true }

func TestDigestSchemaFileNotFound(t *testing.T) {
	_, err := schema.DigestSchema("/path/to/unknown/file.yml")
	assert.ErrorIs(t, err, schema.ErrSchemaFile)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestDigestSchemaUnmarshalError(t *testing.T) {
	_, err := schema.DigestSchema("../test/unit/schema_test_unmarshal_error.yml")
	assert.ErrorIs(t, err, schema.ErrSchemaFile)
	assert.Contains(t, err.Error(), "yaml: unmarshal errors")
}

func TestDigestSchemaUnmarshalErrorUnknownField(t *testing.T) {
	_, err := schema.DigestSchema("../test/unit/schema_test_unmarshal_error_unknown_field.yml")
	assert.Contains(t, err.Error(), "line 7: field ArbitraryField not found in type schema.Table")
}

func TestDigestSchemaValidationError(t *testing.T) {
	_, err := schema.DigestSchema("../test/unit/schema_test.yml")
	assert.ErrorIs(t, err, schema.ErrSchemaValidation)
}

func TestDigestSchema(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	dataconv.RegisterConverter("conv_timestamp", DummyConverter(""))
	schema, err := schema.DigestSchema("../test/unit/schema_test.yml")

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
func TestDigestSchemaReferences(t *testing.T) {
	schema, err := schema.DigestSchema("../test/unit/schema_test_grouping.yml")
	require.Nil(t, err)

	keys := make([]string, 0, len(schema.Refs))
	for k := range schema.Refs {
		keys = append(keys, k)
	}

	assert.Len(t, keys, 8)
	assert.Contains(t, schema.Refs, "b1.id")
	assert.Contains(t, schema.Refs, "a11.id")
	assert.Contains(t, schema.Refs, "b1312.id")
	assert.Contains(t, schema.Refs, "a1.id")
	assert.Contains(t, schema.Refs, "b13.id")
	assert.Contains(t, schema.Refs, "b131.id")
	assert.Contains(t, schema.Refs, "c1.id")
	assert.Contains(t, schema.Refs, "a11.id[@]")
}

func TestSelectColumns(t *testing.T) {
	table := schema.Table{
		Columns: []schema.Column{"id", "name", "description"},
	}

	columns := table.SelectColumns()
	assert.Len(t, columns, 3)
	assert.EqualValues(t, "id", columns[0])
	assert.EqualValues(t, "name", columns[1])
	assert.EqualValues(t, "description", columns[2])

	table = schema.Table{
		Ignore: []schema.Ignore{"tax", "total"},
	}

	columns = table.SelectColumns()
	assert.Len(t, columns, 2)
	assert.EqualValues(t, "tax", columns[0])
	assert.EqualValues(t, "total", columns[1])

	table = schema.Table{}

	columns = table.SelectColumns()
	assert.Len(t, columns, 1)
	assert.EqualValues(t, "*", columns[0])
}

func TestFormattedSelectColumns(t *testing.T) {
	table := schema.Table{
		Columns: []schema.Column{"id", "name", "description"},
	}

	assert.Equal(t, "'id', 'name', 'description'", table.FormattedSelectColumns())

	table = schema.Table{
		Ignore: []schema.Ignore{"tax", "total"},
	}

	assert.Equal(t, "'tax', 'total'", table.FormattedSelectColumns())
}
