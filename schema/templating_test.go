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
  - name: test`

	schema, err := schema.ApplyTemplates("", text)
	assert.Nil(t, err)
	assert.Equal(t, text, schema)
}

func TestApplyTemplatesErrorInvalidTemplateDefinition(t *testing.T) {
	schemaPath := "../test/unit/templating_test.yml"
	text := `---
tables:
  - name: test
  <%= template path="_domain-customer.yml" param 123 %>`

	_, err := schema.ApplyTemplates(schemaPath, text)
	assert.Equal(t,
		"invalid template definition `<%= template path=\"_domain-customer.yml\" param 123 %>'", err.Error())
}

func TestApplyTemplatesErrorEmptyParameter(t *testing.T) {
	schemaPath := "../test/unit/templating_test.yml"
	text := `---
tables:
  - name: test
  <%= template path="_domain-customer.yml" param="" %>`

	_, err := schema.ApplyTemplates(schemaPath, text)
	assert.Equal(t,
		"template parameter 'param' is empty", err.Error())
}

func TestApplyTemplatesErrorRepeatedParameter(t *testing.T) {
	schemaPath := "../test/unit/templating_test.yml"
	text := `---
tables:
  - name: test
  <%= template path="_domain-customer.yml" param="123" param="321" %>`

	_, err := schema.ApplyTemplates(schemaPath, text)
	assert.Equal(t,
		"repeated parameter 'param'", err.Error())
}

func TestApplyTemplatesErrorPathIsRequired(t *testing.T) {
	schemaPath := "../test/unit/templating_test.yml"
	text := `---
tables:
  - name: test
  <%= template param="123" %>`

	_, err := schema.ApplyTemplates(schemaPath, text)
	assert.Equal(t,
		"path parameter is required `<%= template param=\"123\" %>'", err.Error())
}

func TestApplyTemplatesErrorPathNotFound(t *testing.T) {
	schemaPath := "../test/unit/templating_test.yml"
	text := `---
tables:
  - name: test
  <%= template path="/path/to/nowhere" param="123" %>`

	_, err := schema.ApplyTemplates(schemaPath, text)
	assert.Equal(t,
		"/path/to/nowhere not found", err.Error())
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
