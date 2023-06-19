package schema_test

import (
	"os"
	"testing"

	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyTemplatesNoTemplateToRender(t *testing.T) {
	text := `---
tables:
  -name: test`

	schema, err := schema.ApplyTemplates("", text)
	assert.Nil(t, err)
	assert.Equal(t, text, schema)
}

func TestApplyTemplates(t *testing.T) {
	schemaPath := "../test/unit/templating_test.yml"
	bytes, err := os.ReadFile(schemaPath)
	require.Nil(t, err)

	text := string(bytes)
	schema, err := schema.ApplyTemplates(schemaPath, text)
	expected := `---
tables:
  - name: customers
    filters:
      - name: id
        value: ${customer_id}
  - name: addresses
    filters:
      - name: id
        value: ${customer_id}
  - name: preferences
    filters:
      - name: id
        value: ${customer_id}
  - name: orders
    filters:
      - name: customer_id
        value: ${customer_id}
  - name: reviews
    filters:
      - name: id
        value: 123
  - name: categories
    filters:
      - name: id
        value: 123
  - name: product
    filters:
      - name: id
        value: 123`
	assert.Nil(t, err)
	assert.Equal(t, expected, schema)
}
