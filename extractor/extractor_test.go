package extractor_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
)

func TestExtractSchemaFileNotFound(t *testing.T) {
	err := extractor.Extract(extractor.Conf{SchemaPath: ""}, nil)
	assert.ErrorIs(t, err, schema.ErrSchemaFile)
}

func TestExtractUnsupportedReader(t *testing.T) {
	err := extractor.Extract(extractor.Conf{SchemaPath: "../test/unit/schema_test_grouping.yml"}, nil)
	assert.ErrorIs(t, err, reader.ErrUnsupportedDBReader)
}
