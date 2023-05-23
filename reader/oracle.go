package reader

import "github.com/aureliano/db-unit-extractor/schema"

type OracleReader struct {
}

func (r OracleReader) FetchColumnsMetadata(_ schema.Table) ([]DBColumn, error) {
	panic("not implemented yet")
}

func (r OracleReader) FetchData(_ string, _ []DBColumn, _ []string,
	filters [][]interface{}) ([]map[string]interface{}, error) {
	panic("not implemented yet")
}
