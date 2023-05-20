package extractor_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassifyGroupOneNotClassified(t *testing.T) {
	s := extractor.Schema{
		Tables: []extractor.TableSchema{
			{Name: "t1", Filters: []extractor.FilterSchema{{Name: "id", Value: "${table.column}"}}},
		},
	}

	err := s.Classify()
	assert.ErrorIs(t, err, extractor.ErrTableClassification)
	assert.Contains(t, err.Error(), "couldn't find any level one tables")
}

func TestClassifyReferenceNotFound(t *testing.T) {
	s := extractor.Schema{
		Tables: []extractor.TableSchema{
			{Name: "t1", Filters: []extractor.FilterSchema{{Name: "id", Value: "1"}}},
			{Name: "t2", Filters: []extractor.FilterSchema{{Name: "id", Value: "${table.column}"}}},
		},
	}

	err := s.Classify()
	assert.ErrorIs(t, err, extractor.ErrTableClassification)
	assert.Contains(t, err.Error(), "t2.id points to unresolvable reference '${table.column}'")
}

func TestClassify(t *testing.T) {
	schema, err := extractor.DigestSchema("../test/unit/schema_test_grouping.yml")
	require.Nil(t, err)

	err = schema.Classify()
	assert.Nil(t, err)

	assert.Equal(t, 1, schema.Tables[1].GroupID)
	assert.Equal(t, 1, schema.Tables[3].GroupID)
	assert.Equal(t, 1, schema.Tables[7].GroupID)

	assert.Equal(t, 2, schema.Tables[0].GroupID)
	assert.Equal(t, 2, schema.Tables[5].GroupID)
	assert.Equal(t, 2, schema.Tables[6].GroupID)
	assert.Equal(t, 2, schema.Tables[8].GroupID)

	assert.Equal(t, 3, schema.Tables[2].GroupID)
	assert.Equal(t, 3, schema.Tables[9].GroupID)
	assert.Equal(t, 3, schema.Tables[12].GroupID)

	assert.Equal(t, 4, schema.Tables[10].GroupID)
	assert.Equal(t, 4, schema.Tables[11].GroupID)

	assert.Equal(t, 5, schema.Tables[4].GroupID)
}
