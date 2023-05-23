package reader_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
)

func TestFetchColumnsMetadata(t *testing.T) {
	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"})
	assert.Panics(t, func() { r.FetchColumnsMetadata(schema.Table{}) }, "not implemented yet")
}

func TestFetchData(t *testing.T) {
	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"})
	assert.Panics(
		t,
		func() { r.FetchData("", []reader.DBColumn{}, []string{}, [][]interface{}{}) },
		"not implemented yet",
	)
}
