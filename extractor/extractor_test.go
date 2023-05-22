package extractor_test

import (
	"fmt"
	"testing"

	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
)

type DummyReader struct{}

func (DummyReader) FetchColumnsMetadata(table schema.Table) ([]reader.DBColumn, error) {
	return []reader.DBColumn{
		{Name: "id", Type: "int"},
		{Name: "name", Type: "varchar"},
		{Name: "description", Type: "varchar"},
	}, nil
}

func (DummyReader) FetchData(table string, fields []reader.DBColumn, converters []string,
	filters [][]interface{}) ([]map[string]interface{}, error) {
	return nil, nil
}

func TestExtractSchemaFileNotFound(t *testing.T) {
	err := extractor.Extract(extractor.Conf{SchemaPath: ""}, nil)
	assert.ErrorIs(t, err, schema.ErrSchemaFile)
}

func TestExtractUnsupportedReader(t *testing.T) {
	err := extractor.Extract(extractor.Conf{SchemaPath: "../test/unit/schema_test_grouping.yml"}, nil)
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}

func TestExtract(t *testing.T) {
	refs := make(map[string]interface{})
	refs["b1_id"] = 34

	err := extractor.Extract(
		extractor.Conf{
			SchemaPath: "../test/unit/schema_test_grouping.yml",
			References: refs,
		}, DummyReader{},
	)
	fmt.Println(err)
}
