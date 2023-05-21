package reader

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

type DBReader interface {
	FetchColumnsMetadata(table string, fieldsIn, fieldsOut []string) ([]DBColumn, error)
	FetchData(table string, fields []DBColumn, converters []string,
		filters [][]interface{}) ([]map[string]interface{}, error)
}
