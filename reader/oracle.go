package reader

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/schema"
)

type OracleReader struct {
	db        *sql.DB
	profiling *os.File
}

func (r OracleReader) FetchColumnsMetadata(table schema.Table) ([]DBColumn, error) {
	query := buildOracleSQLQueryColumnsMetadata(table)

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Oracle.FetchColumnsMetadata\nTable: %s\nQuery: %s\nError: %s\n", table.Name, query, err.Error())
		return nil, err
	}
	defer rows.Close()

	records := make([]DBColumn, 0, len(table.SelectColumns()))

	for rows.Next() {
		rec := DBColumn{}
		var nullable string
		var precision, scale sql.NullInt64

		err = rows.Scan(&rec.Name, &rec.Type, &nullable, &rec.Length, &precision, &scale)
		if err != nil {
			log.Printf("Oracle.FetchColumnsMetadata\nTable: %s\nScan error: %s\n", table.Name, err.Error())
			return nil, err
		}

		rec.DecimalSize.Precision = precision.Int64
		rec.DecimalSize.Scale = scale.Int64

		rec.Nullable = strToBool(nullable)
		records = append(records, rec)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Oracle.FetchColumnsMetadata\nTable: %s\nFetch rows error: %s\n", table.Name, err.Error())
		return nil, err
	} else if len(records) == 0 {
		return nil, fmt.Errorf(
			"no metadata found for table %s (make sure it exists and user has proper grants)", table.Name)
	}

	return records, nil
}

func (r OracleReader) FetchData(table string, fields []DBColumn, converters []dataconv.Converter,
	filters [][]interface{}) ([][]*DBColumn, error) {
	query := buildOracleSQLQueryColumns(table, fields, filters)

	size := len(filters)
	ind := make([]int, 0)
	values := make([]interface{}, size)
	expectedArrValuesSize := 0

	for i := 0; i < size; i++ {
		value, multivalued := filters[i][1].([]interface{})
		if multivalued {
			ind = append(ind, i)
			expectedArrValuesSize += len(value)
		} else {
			values[i] = filters[i][1]
		}
	}

	arrValues := make([][]interface{}, 0)
	if len(ind) == 0 && !emptyFilter(values) {
		arrValues = append(arrValues, values)
		expectedArrValuesSize++
	}

	for _, i := range ind {
		for _, v := range filters[i][1].([]interface{}) {
			cpValues := make([]interface{}, len(values))
			copy(cpValues, values)
			cpValues[i] = v

			if !emptyFilter(cpValues) {
				arrValues = append(arrValues, cpValues)
			}
		}
	}

	wrongValuesSize := len(arrValues) != expectedArrValuesSize
	filterWithoutValues := len(arrValues) == 0 && len(filters) > 0
	if wrongValuesSize || filterWithoutValues {
		return nil, fmt.Errorf("not all filters were bound for table %s `%v'", table, filters)
	}

	return fetchData(r.db, fields, converters, arrValues, query)
}

func (r OracleReader) ProfilerMode() bool {
	return r.profiling != nil
}

func (r OracleReader) StartDBProfiler(ctx context.Context) {
	go func(c context.Context) {
		for {
			select {
			case <-c.Done():
				_ = r.profiling.Close()
				return
			default:
				r.profileSnapshot()
				time.Sleep(DBSnapshotDealy)
			}
		}
	}(ctx)
}

func fetchData(db *sql.DB, fields []DBColumn, converters []dataconv.Converter, arrValues [][]interface{},
	query string) ([][]*DBColumn, error) {
	rows := make([][]*DBColumn, 0, len(arrValues))

	if len(arrValues) > 0 {
		for _, filterValues := range arrValues {
			data, err := executeQuery(db, fields, converters, filterValues, query)
			if err != nil {
				return nil, err
			}
			rows = append(rows, data...)
		}
	} else {
		data, err := executeQuery(db, fields, converters, []interface{}{}, query)
		if err != nil {
			return nil, err
		}
		rows = append(rows, data...)
	}

	return rows, nil
}

func executeQuery(db *sql.DB, fields []DBColumn, converters []dataconv.Converter, filters []interface{},
	query string) ([][]*DBColumn, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Oracle.executeQuery\nQuery: %s\nPrepare statement error: %s\n", query, err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(filters...)
	if err != nil {
		log.Printf("Oracle.executeQuery\nQuery: %s\nQuery error: %s\n", query, err.Error())
		return nil, err
	}
	defer rows.Close()

	return readDataSet(fields, rows, converters)
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

	if len(table.Columns) == 0 && len(table.Ignore) == 0 {
		return builder.String()
	}

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

func readDataSet(fields []DBColumn, rows *sql.Rows, converters []dataconv.Converter) ([][]*DBColumn, error) {
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	data := make([][]*DBColumn, 0)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		_ = rows.Scan(valuePtrs...)

		row := make([]*DBColumn, count)
		for i, field := range fields {
			f := field
			row[i] = &f
		}

		for i := range columns {
			value, err := fetchValue(values[i], converters)
			if err != nil {
				log.Printf("Oracle.readDataSet\nFetch value error: %s\n", err.Error())
				log.Printf("Field: %s - Value: %v\n", columns[i], values[i])
				return nil, err
			}

			row[i].Value = value
		}

		data = append(data, row)
	}

	return data, rows.Err()
}

func fetchValue(value interface{}, converters []dataconv.Converter) (interface{}, error) {
	converter := findConverter(value, converters)
	if converter != nil {
		return converter.Convert(value)
	}

	return value, nil
}

func findConverter(value interface{}, converters []dataconv.Converter) dataconv.Converter {
	for _, converter := range converters {
		if converter.Handle(value) {
			return converter
		}
	}

	return nil
}

func emptyFilter(filter []interface{}) bool {
	if len(filter) == 0 {
		return true
	}

	for _, f := range filter {
		if f == nil {
			return true
		}
	}

	return false
}

func (r OracleReader) profileSnapshot() {
	if !r.ProfilerMode() {
		return
	}

	stats := r.db.Stats()
	record := fmt.Sprintf("%d %d %d %d %d %d %d %d %d\n", stats.MaxOpenConnections, stats.OpenConnections,
		stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration, stats.MaxIdleClosed, stats.MaxIdleTimeClosed,
		stats.MaxLifetimeClosed)

	_, _ = r.profiling.WriteString(record)
}
