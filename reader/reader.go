package reader

import "fmt"

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

type DBReader interface {
	FetchColumnsMetadata(table string, fieldsIn, fieldsOut []string) ([]DBColumn, error)
	FetchData(table string, fields []DBColumn, converters []string,
		filters [][]interface{}) ([]map[string]interface{}, error)
}

func NewReader(_ DataSource) (DBReader, error) {
	return nil, fmt.Errorf("not implemented yet")
}
