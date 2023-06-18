package schema_test

import (
	"os"
	"testing"

	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyTemplates(t *testing.T) {
	bytes, err := os.ReadFile("../test/unit/templating_test.yml")
	require.Nil(t, err)

	text := string(bytes)
	schema, err := schema.ApplyTemplates(text)
	//assert.Nil(t, err)
	assert.Empty(t, schema)
}
