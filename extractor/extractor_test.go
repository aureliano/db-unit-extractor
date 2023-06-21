package extractor_test

import (
	"context"
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

func (DummyConverter) Convert(interface{}) (interface{}, error) {
	return 0, nil
}

func (DummyConverter) Handle(interface{}) bool { return true }

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

func (DummyReader) FetchData(table string, _ []reader.DBColumn, _ []dataconv.Converter,
	_ [][]interface{}) ([][]*reader.DBColumn, error) {
	switch {
	case table == "customers":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 34},
			{Name: "first_name", Value: "Antonio"},
			{Name: "last_name", Value: "Vivaldi"},
		}}, nil
	case table == "orders":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 63},
			{Name: "status", Value: "paid"},
			{Name: "total", Value: 165.88},
			{Name: "tax", Value: 15.08},
		}}, nil
	case table == "orders_customers":
		return [][]*reader.DBColumn{{
			{Name: "order_id", Value: 63},
			{Name: "customer_id", Value: 34},
		}}, nil
	case table == "products":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 3},
			{Name: "name", Value: "Holy Bible"},
			{Name: "description", Value: "Latin Vulgata from Saint Jerome"},
			{Name: "price", Value: 150.8},
		}}, nil
	default:
		return [][]*reader.DBColumn{}, nil
	}
}

func (HumanResourcesReader) FetchData(table string, _ []reader.DBColumn, _ []dataconv.Converter,
	_ [][]interface{}) ([][]*reader.DBColumn, error) {
	switch {
	case table == "employees":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 100},
			{Name: "department_id", Value: 90},
			{Name: "job_id", Value: 5},
			{Name: "first_name", Value: "Antonio"},
			{Name: "last_name", Value: "Vivaldi"},
		}, {
			{Name: "id", Value: 101},
			{Name: "department_id", Value: 90},
			{Name: "job_id", Value: 3},
			{Name: "first_name", Value: "Johann"},
			{Name: "last_name", Value: "Bach"},
		}}, nil
	case table == "departments":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 90},
			{Name: "location_id", Value: 1700},
			{Name: "name", Value: "Sales"},
		}}, nil
	case table == "locations":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 1700},
			{Name: "country_id", Value: 55},
			{Name: "city", Value: "Governador Valadares"},
			{Name: "province", Value: "Minas Gerais"},
		}}, nil
	case table == "countries":
		return [][]*reader.DBColumn{{{Name: "id", Value: 55},
			{Name: "region_id", Value: 3},
			{Name: "name", Value: "Brasil"},
		}}, nil
	case table == "regions":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 3},
			{Name: "name", Value: "América do Sul"},
		}}, nil
	case table == "jobs":
		return [][]*reader.DBColumn{{
			{Name: "id", Value: 1},
			{Name: "title", Value: "junior developer"},
		}, {
			{Name: "id", Value: 2},
			{Name: "title", Value: "full developer"},
		}, {
			{Name: "id", Value: 3},
			{Name: "title", Value: "senior developer"},
		}, {
			{Name: "id", Value: 4},
			{Name: "title", Value: "seller"},
		}, {
			{Name: "id", Value: 5},
			{Name: "title", Value: "sell manager"},
		}}, nil
	case table == "job_history":
		return [][]*reader.DBColumn{{
			{Name: "employee_id", Value: 100},
			{Name: "job_id", Value: 4},
		}, {
			{Name: "employee_id", Value: 100},
			{Name: "job_id", Value: 5},
		}, {
			{Name: "employee_id", Value: 101},
			{Name: "job_id", Value: 1},
		}, {
			{Name: "employee_id", Value: 101},
			{Name: "job_id", Value: 2},
		}, {
			{Name: "employee_id", Value: 101},
			{Name: "job_id", Value: 3},
		}}, nil
	default:
		return [][]*reader.DBColumn{}, nil
	}
}

func (FetchMetadataErrorDummyReader) FetchData(string, []reader.DBColumn, []dataconv.Converter,
	[][]interface{}) ([][]*reader.DBColumn, error) {
	return [][]*reader.DBColumn{}, nil
}

func (FetchDataErrorDummyReader) FetchData(string, []reader.DBColumn, []dataconv.Converter,
	[][]interface{}) ([][]*reader.DBColumn, error) {
	return nil, fmt.Errorf("fetch data error")
}

func (DummyReader) ProfilerMode() bool {
	return true
}

func (FetchMetadataErrorDummyReader) ProfilerMode() bool {
	return false
}

func (FetchDataErrorDummyReader) ProfilerMode() bool {
	return false
}

func (HumanResourcesReader) ProfilerMode() bool {
	return false
}

func (DummyReader) StartDBProfiler(context.Context) {
}

func (FetchMetadataErrorDummyReader) StartDBProfiler(context.Context) {
}

func (FetchDataErrorDummyReader) StartDBProfiler(context.Context) {
}

func (HumanResourcesReader) StartDBProfiler(context.Context) {
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

func (WriteDataErrorDummyWriter) Write(string, [][]*reader.DBColumn) error {
	return fmt.Errorf("write data error")
}

func (WriteDataHeaderErrorDummyWriter) Write(string, [][]*reader.DBColumn) error {
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
	dataconv.RegisterConverters()

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
