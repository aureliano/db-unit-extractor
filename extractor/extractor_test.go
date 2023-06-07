package extractor_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/aureliano/db-unit-extractor/writer"
	"github.com/stretchr/testify/assert"
)

type DummyConverter string

func (DummyConverter) Convert(_ interface{}, _ *interface{}) {
}

type DummyReader struct{}

type HumanResourcesReader struct{}

type FetchMetadataErrorDummyReader struct{}

type FetchDataErrorDummyReader struct{}

type DummyWriter struct{}

type WriteDataErrorDummyWriter struct{}

type WriteDataHeaderErrorDummyWriter struct{}

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

func (HumanResourcesReader) FetchColumnsMetadata(table schema.Table) ([]reader.DBColumn, error) {
	switch {
	case table.Name == "employees":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "department_id", Type: "int"},
			{Name: "job_id_id", Type: "int"},
			{Name: "first_name", Type: "varchar"},
			{Name: "last_name", Type: "varchar"},
		}, nil
	case table.Name == "departments":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "localtion_id", Type: "int"},
			{Name: "name", Type: "varchar"},
		}, nil
	case table.Name == "locations":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "country_id", Type: "int"},
			{Name: "city", Type: "varchar"},
			{Name: "province", Type: "varchar"},
		}, nil
	case table.Name == "countries":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "region_id", Type: "int"},
			{Name: "name", Type: "varchar"},
		}, nil
	case table.Name == "regions":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "varchar"},
		}, nil
	case table.Name == "jobs":
		return []reader.DBColumn{
			{Name: "id", Type: "int"},
			{Name: "title", Type: "varchar"},
		}, nil
	case table.Name == "job_history":
		return []reader.DBColumn{
			{Name: "department_id", Type: "int"},
			{Name: "employee_id", Type: "int"},
		}, nil
	default:
		return []reader.DBColumn{}, nil
	}
}

func (FetchMetadataErrorDummyReader) FetchColumnsMetadata(schema.Table) ([]reader.DBColumn, error) {
	return nil, fmt.Errorf("fetch metadata error")
}

func (FetchDataErrorDummyReader) FetchColumnsMetadata(schema.Table) ([]reader.DBColumn, error) {
	return []reader.DBColumn{}, nil
}

func (DummyReader) FetchData(table string, _ []reader.DBColumn, _ []schema.Converter,
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

func (HumanResourcesReader) FetchData(table string, _ []reader.DBColumn, _ []schema.Converter,
	_ [][]interface{}) ([]map[string]interface{}, error) {
	switch {
	case table == "employees":
		m1 := make(map[string]interface{})
		m1["id"] = 100
		m1["department_id"] = 90
		m1["job_id"] = 5
		m1["first_name"] = "Antonio"
		m1["last_name"] = "Vivaldi"

		m2 := make(map[string]interface{})
		m2["id"] = 101
		m2["department_id"] = 90
		m2["job_id"] = 3
		m2["first_name"] = "Johann"
		m2["last_name"] = "Bach"

		return []map[string]interface{}{m1, m2}, nil
	case table == "departments":
		m := make(map[string]interface{})
		m["id"] = 90
		m["location_id"] = 1700
		m["name"] = "Sales"
		return []map[string]interface{}{m}, nil
	case table == "locations":
		m := make(map[string]interface{})
		m["id"] = 1700
		m["country_id"] = 55
		m["city"] = "Governador Valadares"
		m["province"] = "Minas Gerais"
		return []map[string]interface{}{m}, nil
	case table == "countries":
		m := make(map[string]interface{})
		m["id"] = 55
		m["region_id"] = 3
		m["name"] = "Brasil"
		return []map[string]interface{}{m}, nil
	case table == "regions":
		m := make(map[string]interface{})
		m["id"] = 3
		m["name"] = "América do Sul"
		return []map[string]interface{}{m}, nil
	case table == "jobs":
		m1 := make(map[string]interface{})
		m1["id"] = 1
		m1["title"] = "junior developer"

		m2 := make(map[string]interface{})
		m2["id"] = 2
		m2["title"] = "full developer"

		m3 := make(map[string]interface{})
		m3["id"] = 3
		m3["title"] = "senior developer"

		m4 := make(map[string]interface{})
		m4["id"] = 4
		m4["title"] = "seller"

		m5 := make(map[string]interface{})
		m5["id"] = 5
		m5["title"] = "sell manager"

		return []map[string]interface{}{m1, m2, m3, m4, m5}, nil
	case table == "job_history":
		m1 := make(map[string]interface{})
		m1["employee_id"] = 100
		m1["job_id"] = 4

		m2 := make(map[string]interface{})
		m2["employee_id"] = 100
		m2["job_id"] = 5

		m3 := make(map[string]interface{})
		m3["employee_id"] = 101
		m3["job_id"] = 1

		m4 := make(map[string]interface{})
		m4["employee_id"] = 101
		m4["job_id"] = 2

		m5 := make(map[string]interface{})
		m5["employee_id"] = 101
		m5["job_id"] = 3

		return []map[string]interface{}{m1, m2, m3, m4, m5}, nil
	default:
		return []map[string]interface{}{}, nil
	}
}

func (FetchMetadataErrorDummyReader) FetchData(string, []reader.DBColumn, []schema.Converter,
	[][]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (FetchDataErrorDummyReader) FetchData(string, []reader.DBColumn, []schema.Converter,
	[][]interface{}) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("fetch data error")
}

func (DummyWriter) WriteHeader() error { return nil }

func (WriteDataErrorDummyWriter) WriteHeader() error { return nil }

func (WriteDataHeaderErrorDummyWriter) WriteHeader() error { return fmt.Errorf("write header error") }

func (DummyWriter) WriteFooter() error { return nil }

func (WriteDataErrorDummyWriter) WriteFooter() error { return nil }

func (WriteDataHeaderErrorDummyWriter) WriteFooter() error { return nil }

func (DummyWriter) Write(string, []map[string]interface{}) error {
	return nil
}

func (WriteDataErrorDummyWriter) Write(string, []map[string]interface{}) error {
	return fmt.Errorf("write data error")
}

func (WriteDataHeaderErrorDummyWriter) Write(string, []map[string]interface{}) error {
	return nil
}

func TestExtractSchemaFileNotFound(t *testing.T) {
	err := extractor.Extract(extractor.Conf{SchemaPath: ""}, nil, nil)
	assert.ErrorIs(t, err, schema.ErrSchemaFile)
}

func TestExtractUnsupportedReader(t *testing.T) {
	err := extractor.Extract(
		extractor.Conf{SchemaPath: "../test/unit/schema_test_grouping.yml"}, nil, nil,
	)
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}

func TestExtractUnsupportedWriter(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["customer_id"] = 34

	err := extractor.Extract(
		extractor.Conf{
			SchemaPath:  "../test/unit/extractor_test.yml",
			References:  refs,
			OutputTypes: []string{"unknown"},
		}, DummyReader{}, nil,
	)

	assert.ErrorIs(t, err, writer.ErrUnsupportedFileWriter)
}

func TestExtractUnresolvableFilter(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	err := extractor.Extract(
		extractor.Conf{
			SchemaPath: "../test/unit/extractor_test.yml",
		}, DummyReader{}, nil,
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
		}, FetchMetadataErrorDummyReader{}, nil,
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
		}, FetchDataErrorDummyReader{}, nil,
	)

	assert.ErrorIs(t, err, extractor.ErrExtractor)
	assert.Contains(t, err.Error(), "fetch data error")
}

func TestExtractWriteDataError(t *testing.T) {
	mu := sync.Mutex{}
	var handledError error
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		mu.Lock()
		handledError = fmt.Errorf("write data panic")
		mu.Unlock()
	})
	defer patches.Reset()

	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["customer_id"] = 34

	conf := extractor.Conf{
		SchemaPath: "../test/unit/extractor_test.yml",
		References: refs,
	}

	_ = extractor.Extract(conf, DummyReader{}, []writer.FileWriter{WriteDataErrorDummyWriter{}})
	assert.Contains(t, handledError.Error(), "write data panic")
}

func TestExtractWriteDataHeaderError(t *testing.T) {
	mu := sync.Mutex{}
	var handledError error
	patches := gomonkey.ApplyFunc(os.Exit, func(int) {
		mu.Lock()
		handledError = fmt.Errorf("write data panic")
		mu.Unlock()
	})
	defer patches.Reset()

	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["customer_id"] = 34

	conf := extractor.Conf{
		SchemaPath: "../test/unit/extractor_test.yml",
		References: refs,
	}

	_ = extractor.Extract(conf, DummyReader{}, []writer.FileWriter{WriteDataHeaderErrorDummyWriter{}})
	assert.Contains(t, handledError.Error(), "write data panic")
}

func TestExtract(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["customer_id"] = 34

	err := extractor.Extract(
		extractor.Conf{
			SchemaPath:  "../test/unit/extractor_test.yml",
			References:  refs,
			OutputTypes: []string{"console"},
		}, DummyReader{}, nil,
	)

	assert.Nil(t, err)
}

func TestExtractMultivalued(t *testing.T) {
	dataconv.RegisterConverter("conv_date_time", DummyConverter(""))
	refs := make(map[string]interface{})
	refs["department_id"] = 12

	err := extractor.Extract(
		extractor.Conf{
			SchemaPath:  "../test/unit/extractor_multivalued_test.yml",
			References:  refs,
			OutputTypes: []string{"console"},
		}, HumanResourcesReader{}, nil,
	)

	assert.Nil(t, err)
}
