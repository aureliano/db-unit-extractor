package reader

import (
	"database/sql"
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

type DataSource struct {
	DBMSName    string
	Username    string
	Password    string
	Database    string
	Hostname    string
	Port        int
	MaxOpenConn int
	MaxIdleConn int
}

var ErrUnsupportedDBReader = errors.New("unsupported database")

type DBReader interface {
	FetchColumnsMetadata(table schema.Table) ([]DBColumn, error)
	FetchData(table string, fields []DBColumn, converters []schema.Converter,
		filters [][]interface{}) ([]map[string]interface{}, error)
}

func NewReader(ds DataSource, db *sql.DB) (DBReader, error) {
	if strings.ToLower(ds.DBMSName) == "oracle" {
		return OracleReader{db: db}, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedDBReader, ds.DBMSName)
}
