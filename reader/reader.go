package reader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aureliano/db-unit-extractor/schema"
)

type DecimalColumn struct {
	Precision int64
	Scale     int64
}

type DBColumn struct {
	Name        string
	Type        string
	Nullable    bool
	Length      int64
	DecimalSize DecimalColumn
}

var ErrUnsupportedDBReader = errors.New("unsupported database")

type DBReader interface {
	FetchColumnsMetadata(table schema.Table) ([]DBColumn, error)
	FetchData(table string, fields []DBColumn, converters []schema.Converter,
		filters [][]interface{}) ([]map[string]interface{}, error)
}

func NewReader(ds *DataSource) (DBReader, error) {
	if strings.EqualFold(ds.DriverName(), "oracle") {
		return OracleReader{db: ds.DB}, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedDBReader, ds.DriverName())
}
