package reader

import (
	"context"
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
	FetchColumnsMetadata(schema.Table) ([]DBColumn, error)
	FetchData(string, []DBColumn, []dataconv.Converter, [][]interface{}) ([]map[string]interface{}, error)
	ProfilerMode() bool
	StartDBProfiler(context.Context)
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
