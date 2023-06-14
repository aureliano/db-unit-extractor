package reader_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DummyConverter struct{}

func (DummyConverter) Convert(interface{}) (interface{}, error) {
	return nil, fmt.Errorf("converter error")
}

func (DummyConverter) Handle(interface{}) bool {
	return true
}

func TestFetchColumnsMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	rows := sqlmock.NewRows([]string{
		"COLUMN_NAME", "DATA_TYPE", "NULLABLE", "DATA_LENGTH", "DATA_PRECISION", "DATA_SCALE",
	}).
		AddRow("ID", "NUMBER", "F", 22, 2, 0).
		AddRow("USER_ID", "NUMBER", "F", 22, 2, 0).
		AddRow("STATUS", "VARCHAR2", "Y", 15, nil, nil).
		AddRow("TOTAL", "NUMBER", "F", 22, 17, 2)

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnRows(rows)

	columns, err := r.FetchColumnsMetadata(schema.Table{
		Name:    "orders",
		Columns: []schema.Column{"id", "user_id", "status", "total"},
	})

	require.Nil(t, err)

	assert.Equal(t, "ID", columns[0].Name)
	assert.Equal(t, "NUMBER", columns[0].Type)
	assert.Equal(t, false, columns[0].Nullable)
	assert.EqualValues(t, 22, columns[0].Length)
	assert.EqualValues(t, 2, columns[0].DecimalSize.Precision)
	assert.EqualValues(t, 0, columns[0].DecimalSize.Scale)

	assert.Equal(t, "USER_ID", columns[1].Name)
	assert.Equal(t, "NUMBER", columns[1].Type)
	assert.Equal(t, false, columns[1].Nullable)
	assert.EqualValues(t, 22, columns[1].Length)
	assert.EqualValues(t, 2, columns[1].DecimalSize.Precision)
	assert.EqualValues(t, 0, columns[1].DecimalSize.Scale)

	assert.Equal(t, "STATUS", columns[2].Name)
	assert.Equal(t, "VARCHAR2", columns[2].Type)
	assert.Equal(t, true, columns[2].Nullable)
	assert.EqualValues(t, 15, columns[2].Length)
	assert.EqualValues(t, 0, columns[2].DecimalSize.Precision)
	assert.EqualValues(t, 0, columns[2].DecimalSize.Scale)

	assert.Equal(t, "TOTAL", columns[3].Name)
	assert.Equal(t, "NUMBER", columns[3].Type)
	assert.Equal(t, false, columns[3].Nullable)
	assert.EqualValues(t, 22, columns[3].Length)
	assert.EqualValues(t, 17, columns[3].DecimalSize.Precision)
	assert.EqualValues(t, 2, columns[3].DecimalSize.Scale)
}

func TestFetchColumnsMetadataIgnore(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	rows := sqlmock.NewRows([]string{
		"COLUMN_NAME", "DATA_TYPE", "NULLABLE", "DATA_LENGTH", "DATA_PRECISION", "DATA_SCALE",
	}).
		AddRow("ID", "NUMBER", "F", 22, 2, 0).
		AddRow("TOTAL", "NUMBER", "F", 22, 17, 2)

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnRows(rows)

	columns, err := r.FetchColumnsMetadata(schema.Table{
		Name:   "orders",
		Ignore: []schema.Ignore{"city", "status", "user_id"},
	})

	require.Nil(t, err)

	assert.Equal(t, "ID", columns[0].Name)
	assert.Equal(t, "NUMBER", columns[0].Type)
	assert.Equal(t, false, columns[0].Nullable)
	assert.EqualValues(t, 22, columns[0].Length)
	assert.EqualValues(t, 2, columns[0].DecimalSize.Precision)
	assert.EqualValues(t, 0, columns[0].DecimalSize.Scale)

	assert.Equal(t, "TOTAL", columns[1].Name)
	assert.Equal(t, "NUMBER", columns[1].Type)
	assert.Equal(t, false, columns[1].Nullable)
	assert.EqualValues(t, 22, columns[1].Length)
	assert.EqualValues(t, 17, columns[1].DecimalSize.Precision)
	assert.EqualValues(t, 2, columns[1].DecimalSize.Scale)
}

func TestFetchColumnsMetadataAllColumns(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	rows := sqlmock.NewRows([]string{
		"COLUMN_NAME", "DATA_TYPE", "NULLABLE", "DATA_LENGTH", "DATA_PRECISION", "DATA_SCALE",
	}).
		AddRow("ID", "NUMBER", "F", 22, 2, 0).
		AddRow("USER_ID", "NUMBER", "F", 22, 2, 0).
		AddRow("STATUS", "VARCHAR2", "Y", 15, nil, nil).
		AddRow("TOTAL", "NUMBER", "F", 22, 17, 2)

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnRows(rows)

	columns, err := r.FetchColumnsMetadata(schema.Table{Name: "orders"})

	require.Nil(t, err)

	assert.Equal(t, "ID", columns[0].Name)
	assert.Equal(t, "NUMBER", columns[0].Type)
	assert.Equal(t, false, columns[0].Nullable)
	assert.EqualValues(t, 22, columns[0].Length)
	assert.EqualValues(t, 2, columns[0].DecimalSize.Precision)
	assert.EqualValues(t, 0, columns[0].DecimalSize.Scale)

	assert.Equal(t, "USER_ID", columns[1].Name)
	assert.Equal(t, "NUMBER", columns[1].Type)
	assert.Equal(t, false, columns[1].Nullable)
	assert.EqualValues(t, 22, columns[1].Length)
	assert.EqualValues(t, 2, columns[1].DecimalSize.Precision)
	assert.EqualValues(t, 0, columns[1].DecimalSize.Scale)

	assert.Equal(t, "STATUS", columns[2].Name)
	assert.Equal(t, "VARCHAR2", columns[2].Type)
	assert.Equal(t, true, columns[2].Nullable)
	assert.EqualValues(t, 15, columns[2].Length)
	assert.EqualValues(t, 0, columns[2].DecimalSize.Precision)
	assert.EqualValues(t, 0, columns[2].DecimalSize.Scale)

	assert.Equal(t, "TOTAL", columns[3].Name)
	assert.Equal(t, "NUMBER", columns[3].Type)
	assert.Equal(t, false, columns[3].Nullable)
	assert.EqualValues(t, 22, columns[3].Length)
	assert.EqualValues(t, 17, columns[3].DecimalSize.Precision)
	assert.EqualValues(t, 2, columns[3].DecimalSize.Scale)
}

func TestFetchColumnsMetadataQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})
	sqlErr := errors.New("ORA-00923 invalid query")

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnError(sqlErr)

	_, err = r.FetchColumnsMetadata(schema.Table{Name: "customers", Ignore: []schema.Ignore{"id"}})
	assert.ErrorIs(t, err, sqlErr)
}

func TestFetchColumnsMetadataScanRowError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	rows := sqlmock.NewRows([]string{
		"COLUMN_NAME", "DATA_TYPE", "NULLABLE", "DATA_LENGTH", "DATA_PRECISION", "DATA_SCALE",
	}).AddRow("TOTAL", "NUMBER", "F", "22r", 17, 2)

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnRows(rows)

	_, err = r.FetchColumnsMetadata(schema.Table{Name: "customers", Ignore: []schema.Ignore{"id"}})
	assert.Contains(t, err.Error(), "Scan error on column index 3, name \"DATA_LENGTH\"")
}

func TestFetchColumnsMetadataRowsError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})
	sqlErr := errors.New("rows error")

	rows := sqlmock.NewRows([]string{
		"COLUMN_NAME", "DATA_TYPE", "NULLABLE", "DATA_LENGTH", "DATA_PRECISION", "DATA_SCALE",
	}).AddRow("TOTAL", "NUMBER", "F", 22, 17, 2).RowError(0, sqlErr)

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnRows(rows)

	_, err = r.FetchColumnsMetadata(schema.Table{Name: "customers", Ignore: []schema.Ignore{"id"}})
	assert.ErrorIs(t, err, sqlErr)
}

func TestFetchColumnsMetadataEmptyResultError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	rows := sqlmock.NewRows([]string{
		"COLUMN_NAME", "DATA_TYPE", "NULLABLE", "DATA_LENGTH", "DATA_PRECISION", "DATA_SCALE",
	})

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnRows(rows)

	_, err = r.FetchColumnsMetadata(schema.Table{Name: "customers", Ignore: []schema.Ignore{"id"}})
	assert.Equal(t, "no metadata found for table customers (make sure it exists and user has proper grants)", err.Error())
}

func TestFetchData(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	fields := []reader.DBColumn{
		{Name: "ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "USER_ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "STATUS", Type: "VARCHAR2", Nullable: true, Length: 15},
		{Name: "TOTAL", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 17, Scale: 2}},
		{Name: "DATE_REC", Type: "DATE", Nullable: false},
		{Name: "ATTACHMENT", Type: "BLOB", Nullable: true, Length: 1024},
	}
	converters := []dataconv.Converter{dataconv.DateTimeISO8601Converter{}, dataconv.BlobConverter{}}
	filters := [][]interface{}{}
	dateRec := time.Now()
	attachment := []byte("hello world")

	rows := sqlmock.
		NewRows([]string{"ID", "USER_ID", "STATUS", "TOTAL", "DATE_REC", "ATTACHMENT"}).
		AddRow(4, 375, "SOLD", 2243.72, dateRec, attachment)

	sql := "^SELECT (.+) FROM ORDERS$"
	mock.ExpectPrepare(sql).ExpectQuery().WillReturnRows(rows)

	data, err := r.FetchData("ORDERS", fields, converters, filters)
	require.Nil(t, err)

	assert.Len(t, data, 1)
	assert.EqualValues(t, 4, data[0]["ID"])
	assert.EqualValues(t, 375, data[0]["USER_ID"])
	assert.Equal(t, "SOLD", data[0]["STATUS"])
	assert.EqualValues(t, 2243.72, data[0]["TOTAL"])
	assert.Equal(t, dateRec.Format("2006-01-02T15:04:05.999 -0700"), data[0]["DATE_REC"])
	assert.Equal(t, "aGVsbG8gd29ybGQ=", data[0]["ATTACHMENT"])
}

func TestFetchDataNotAllFieldsWereBound(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	fields := []reader.DBColumn{
		{Name: "ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "USER_ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "STATUS", Type: "VARCHAR2", Nullable: true, Length: 15},
		{Name: "TOTAL", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 17, Scale: 2}},
		{Name: "DATE_REC", Type: "DATE", Nullable: false},
		{Name: "ATTACHMENT", Type: "BLOB", Nullable: true, Length: 1024},
	}
	converters := []dataconv.Converter{dataconv.DateTimeISO8601Converter{}, dataconv.BlobConverter{}}
	filters := [][]interface{}{{"ID", nil}}
	dateRec := time.Now()
	attachment := []byte("hello world")

	rows := sqlmock.
		NewRows([]string{"ID", "USER_ID", "STATUS", "TOTAL", "DATE_REC", "ATTACHMENT"}).
		AddRow(4, 375, "SOLD", 2243.72, dateRec, attachment)

	sql := "^SELECT (.+) FROM ORDERS WHERE ID = :1$"
	mock.ExpectPrepare(sql).ExpectQuery().WillReturnRows(rows)

	_, err = r.FetchData("ORDERS", fields, converters, filters)
	assert.NotNil(t, err)
	assert.Equal(t, "not all filters were bound for table ORDERS `[[ID <nil>]]'", err.Error())
}

func TestFetchDataFiltered(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	fields := []reader.DBColumn{
		{Name: "ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "USER_ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "STATUS", Type: "VARCHAR2", Nullable: true, Length: 15},
		{Name: "TOTAL", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 17, Scale: 2}},
	}
	converters := []dataconv.Converter{}
	filters := [][]interface{}{{"ID", 4}, {"STATUS", "SOLD"}}

	rows := sqlmock.
		NewRows([]string{"ID", "USER_ID", "STATUS", "TOTAL"}).
		AddRow(4, 375, "SOLD", 2243.72)

	sql := "^SELECT (.+) FROM ORDERS WHERE ID = :1 AND STATUS = :2$"
	mock.ExpectPrepare(sql).ExpectQuery().WillReturnRows(rows)

	data, err := r.FetchData("ORDERS", fields, converters, filters)
	require.Nil(t, err)

	assert.Len(t, data, 1)
	assert.EqualValues(t, 4, data[0]["ID"])
	assert.EqualValues(t, 375, data[0]["USER_ID"])
	assert.Equal(t, "SOLD", data[0]["STATUS"])
	assert.EqualValues(t, 2243.72, data[0]["TOTAL"])
}

func TestFetchDataPrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	fields := []reader.DBColumn{
		{Name: "ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "USER_ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "STATUS", Type: "VARCHAR2", Nullable: true, Length: 15},
		{Name: "TOTAL", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 19, Scale: 2}},
	}
	converters := []dataconv.Converter{}
	filters := [][]interface{}{}

	errTest := errors.New("prepare error")
	sql := "^SELECT (.+) FROM ORDERS$"
	mock.ExpectPrepare(sql).WillReturnError(errTest)

	_, err = r.FetchData("ORDERS", fields, converters, filters)
	assert.ErrorIs(t, err, errTest)
}

func TestFetchDataPrepareErrorMultivalued(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	fields := []reader.DBColumn{
		{Name: "ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "USER_ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "STATUS", Type: "VARCHAR2", Nullable: true, Length: 15},
		{Name: "TOTAL", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 19, Scale: 2}},
	}
	converters := []dataconv.Converter{}
	filters := [][]interface{}{{"f1", []interface{}{"v1", "v2"}}}

	errTest := errors.New("prepare error")
	sql := "^SELECT (.+) FROM ORDERS WHERE f1 = :1$"
	mock.ExpectPrepare(sql).WillReturnError(errTest)

	_, err = r.FetchData("ORDERS", fields, converters, filters)
	assert.ErrorIs(t, err, errTest)
}

func TestFetchDataPrepareQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	fields := []reader.DBColumn{
		{Name: "ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "USER_ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "STATUS", Type: "VARCHAR2", Nullable: true, Length: 15},
		{Name: "TOTAL", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 18, Scale: 2}},
	}
	converters := []dataconv.Converter{}
	filters := [][]interface{}{}

	errTest := errors.New("prepare query error")
	sql := "^SELECT (.+) FROM ORDERS$"
	mock.ExpectPrepare(sql).ExpectQuery().WillReturnError(errTest)

	_, err = r.FetchData("ORDERS", fields, converters, filters)
	assert.ErrorIs(t, err, errTest)
}

func TestFetchDataConverterError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(&reader.DataSource{DSN: "oracle://usr:pwd@localhost:1521/dbname", DB: db})

	fields := []reader.DBColumn{
		{Name: "ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "USER_ID", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 2, Scale: 0}},
		{Name: "STATUS", Type: "VARCHAR2", Nullable: true, Length: 15},
		{Name: "TOTAL", Type: "NUMBER", Nullable: false, Length: 22,
			DecimalSize: reader.DecimalColumn{Precision: 17, Scale: 2}},
	}
	converters := []dataconv.Converter{DummyConverter{}}
	filters := [][]interface{}{}

	rows := sqlmock.
		NewRows([]string{"ID", "USER_ID", "STATUS", "TOTAL"}).
		AddRow(4, 375, "SOLD", 2243.72)

	sql := "^SELECT (.+) FROM ORDERS$"
	mock.ExpectPrepare(sql).ExpectQuery().WillReturnRows(rows)

	_, err = r.FetchData("ORDERS", fields, converters, filters)
	assert.Equal(t, "converter error", err.Error())
}
