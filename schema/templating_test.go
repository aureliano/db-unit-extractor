package schema_test

import (
	"os"
	"testing"

	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyTemplates(t *testing.T) {
	schemaPath := "../test/unit/templating_test.yml"
	bytes, err := os.ReadFile(schemaPath)
	require.Nil(t, err)

	text := string(bytes)
	schema, err := schema.ApplyTemplates(schemaPath, text)
	assert.Nil(t, err)
	assert.NotEmpty(t, schema)
}
