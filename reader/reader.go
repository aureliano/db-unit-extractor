package reader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aureliano/db-unit-extractor/dataconv"
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
	FetchData(table string, fields []DBColumn, converters []dataconv.Converter,
		filters [][]interface{}) ([]map[string]interface{}, error)
}

func NewReader(ds DBConnector) (DBReader, error) {
	if strings.EqualFold(ds.DriverName(), "oracle") {
		return newOracle(ds)
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedDBReader, ds.DriverName())
}

func newOracle(ds DBConnector) (DBReader, error) {
	db, err := ds.Connect(MaxDBTimeout)
	return OracleReader{db: db}, err
}
