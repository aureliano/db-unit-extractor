package reader

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aureliano/db-unit-extractor/schema"
)

type OracleReader struct {
	ds *DataSource
}

func (r OracleReader) FetchColumnsMetadata(table schema.Table) ([]DBColumn, error) {
	query := buildOracleSQLQueryColumnsMetadata(table)

	rows, err := r.ds.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]DBColumn, len(table.SelectColumns()))
	i := 0

	for rows.Next() {
		rec := DBColumn{}
		var nullable string
		var precision, scale sql.NullInt64

		err = rows.Scan(&rec.Name, &rec.Type, &nullable, &rec.Length, &precision, &scale)
		if err != nil {
			return nil, err
		}

		rec.DecimalSize.Precision = precision.Int64
		rec.DecimalSize.Scale = scale.Int64

		rec.Nullable = strToBool(nullable)
		records[i] = rec
		i++
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return records, nil
}

func (r OracleReader) FetchData(table string, fields []DBColumn, _ []schema.Converter,
	filters [][]interface{}) ([]map[string]interface{}, error) {
	query := buildOracleSQLQueryColumns(table, fields, filters)

	stmt, err := r.ds.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	size := len(filters)
	values := make([]interface{}, size)
	for i := 0; i < size; i++ {
		values[i] = (filters)[i][1]
	}

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readDataSet(rows)
}

func strToBool(str string) bool {
	return str == "Y"
}

func buildOracleSQLQueryColumnsMetadata(table schema.Table) string {
	var builder strings.Builder
	builder.WriteString("SELECT COLUMN_NAME, DATA_TYPE, NULLABLE, DATA_LENGTH, DATA_PRECISION, DATA_SCALE")
	builder.WriteString(" FROM ALL_TAB_COLS WHERE")
	builder.WriteString(fmt.Sprintf(" TABLE_NAME = '%s'", strings.ToUpper(table.Name)))
	builder.WriteString(" AND VIRTUAL_COLUMN = 'NO'")
	builder.WriteString(" AND COLUMN_NAME NOT LIKE '%$%'")
	builder.WriteString(" AND COLUMN_NAME")

	if len(table.Ignore) > 0 {
		builder.WriteString(" NOT")
	}

	builder.WriteString(fmt.Sprintf(" IN(%s)", strings.ToUpper(table.FormattedSelectColumns())))

	return builder.String()
}

func buildOracleSQLQueryColumns(table string, fields []DBColumn, filters [][]interface{}) string {
	var sql strings.Builder

	sql.WriteString("SELECT ")
	size := len(fields)
	for i := 0; i < size; i++ {
		field := fields[i]

		if field.Type == "NUMBER" && field.DecimalSize.Precision > 18 {
			sql.WriteString(fmt.Sprintf("TO_CHAR(%s) AS %s", field.Name, field.Name))
		} else {
			sql.WriteString(field.Name)
		}

		if i < size-1 {
			sql.WriteString(",")
		}
	}

	sql.WriteString(fmt.Sprintf(" FROM %s", table))

	size = len(filters)
	if size > 0 {
		sql.WriteString(" WHERE ")
	}

	for i := 0; i < size; i++ {
		sql.WriteString(fmt.Sprintf("%s = :%d", (filters)[i][0], i+1))
		if i < size-1 {
			sql.WriteString(" AND ")
		}
	}

	query := sql.String()
	return query
}

func readDataSet(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	data := make([]map[string]interface{}, 0)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		_ = rows.Scan(valuePtrs...)

		row := make(map[string]interface{})
		for i := range columns {
			row[columns[i]] = values[i]
		}

		data = append(data, row)
	}

	return data, rows.Err()
}
