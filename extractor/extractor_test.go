package extractor_test

import (
	"fmt"
	"testing"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
)

type DummyConverter string

func (DummyConverter) Convert(_ interface{}, _ *interface{}) {
}

type DummyReader struct{}

type FetchMetadataErrorDummyReader struct{}

type FetchDataErrorDummyReader struct{}

func (DummyReader) FetchColumnsMetadata(table schema.Table) ([]reader.DBColumn, error) {
	switch {
	case table.Name == "customers":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "first_name", Type: "varchar"},
			{Name: "last_name", Type: "varchar"},
		}, nil
	case table.Name == "orders":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "status", Type: "varchar"},
			{Name: "total", Type: "float"},
			{Name: "tax", Type: "float"},
		}, nil
	case table.Name == "orders_customers":
		return []reader.DBColumn{
			{Name: "order_id", Type: "int"},
			{Name: "customer_id", Type: "int"},
		}, nil
	case table.Name == "products":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "varchar"},
			{Name: "description", Type: "varchar"},
			{Name: "price", Type: "float"},
		}, nil
	default:
		return []reader.DBColumn{}, nil
	}
}

func (FetchMetadataErrorDummyReader) FetchColumnsMetadata(_ schema.Table) ([]reader.DBColumn, error) {
	return nil, fmt.Errorf("fetch metadata error")
}

func (FetchDataErrorDummyReader) FetchColumnsMetadata(_ schema.Table) ([]reader.DBColumn, error) {
	return []reader.DBColumn{}, nil
}

func (DummyReader) FetchData(table string, _ []reader.DBColumn, _ []string,
	_ [][]interface{}) ([]map[string]interface{}, error) {
	switch {
	case table == "customers":
		m := make(map[string]interface{})
		m["id"] = 34
		m["first_name"] = "Antonio"
		m["last_name"] = "Vivaldi"
		return []map[string]interface{}{m}, nil
	case table == "orders":
		m := make(map[string]interface{})
		m["id"] = 63
		m["status"] = "paid"
		m["total"] = 165.88
		m["tax"] = 15.08
		return []map[string]interface{}{m}, nil
	case table == "orders_customers":
		m := make(map[string]interface{})
		m["order_id"] = 63
		m["customer_id"] = 34
		return []map[string]interface{}{m}, nil
	case table == "products":
		m := make(map[string]interface{})
		m["id"] = 3
		m["name"] = "Holy Bible"
		m["description"] = "Latin Vulgata from Saint Jerome"
		m["price"] = 150.8
		return []map[string]interface{}{m}, nil
	default:
		return []map[string]interface{}{}, nil
	}
}

func (FetchMetadataErrorDummyReader) FetchData(_ string, _ []reader.DBColumn, _ []string,
	_ [][]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (FetchDataErrorDummyReader) FetchData(_ string, _ []reader.DBColumn, _ []string,
	_ [][]interface{}) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("fetch data error")
}

func TestExtractSchemaFileNotFound(t *testing.T) {
	err := extractor.Extract(extractor.Conf{SchemaPath: ""}, nil)
	assert.ErrorIs(t, err, schema.ErrSchemaFile)
}

func TestExtractUnsupportedReader(t *testing.T) {
	err := extractor.Extract(extractor.Conf{SchemaPath: "../test/unit/schema_test_grouping.yml"}, nil)
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}

func TestExtractUnresolvableFilter(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	err := extractor.Extract(
		extractor.Conf{
			SchemaPath: "../test/unit/extractor_test.yml",
		}, DummyReader{},
	)

	assert.ErrorIs(t, err, extractor.ErrExtractor)
	assert.Contains(t, err.Error(), "filter customers.id not found '${customer_id}")
}

func TestExtractFetchMetadataError(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["customer_id"] = 34

	err := extractor.Extract(
		extractor.Conf{
			SchemaPath: "../test/unit/extractor_test.yml",
			References: refs,
		}, FetchMetadataErrorDummyReader{},
	)

	assert.ErrorIs(t, err, extractor.ErrExtractor)
	assert.Contains(t, err.Error(), "fetch metadata error")
}

func TestExtractFetchDataError(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["customer_id"] = 34

	err := extractor.Extract(
		extractor.Conf{
			SchemaPath: "../test/unit/extractor_test.yml",
			References: refs,
		}, FetchDataErrorDummyReader{},
	)

	assert.ErrorIs(t, err, extractor.ErrExtractor)
	assert.Contains(t, err.Error(), "fetch data error")
}

func TestExtract(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["customer_id"] = 34

	err := extractor.Extract(
		extractor.Conf{
			SchemaPath: "../test/unit/extractor_test.yml",
			References: refs,
		}, DummyReader{},
	)

	assert.Nil(t, err)
}
