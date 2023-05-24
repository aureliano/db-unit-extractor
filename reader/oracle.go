package reader

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aureliano/db-unit-extractor/schema"
)

type OracleReader struct {
	db *sql.DB
}

func (r OracleReader) FetchColumnsMetadata(table schema.Table) ([]DBColumn, error) {
	query := buildOracleSQLQueryColumnsMetadata(table)

	rows, err := r.db.Query(query)
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

func (r OracleReader) FetchData(_ string, _ []DBColumn, _ []string,
	_ [][]interface{}) ([]map[string]interface{}, error) {
	panic("not implemented yet")
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
