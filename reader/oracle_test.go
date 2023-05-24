package reader_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchColumnsMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"}, db)

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

	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"}, db)

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

func TestFetchColumnsMetadataQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"}, db)
	sqlErr := errors.New("ORA-00923 invalid query")

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnError(sqlErr)

	_, err = r.FetchColumnsMetadata(schema.Table{Name: "customers", Ignore: []schema.Ignore{"id"}})
	assert.ErrorIs(t, err, sqlErr)
}

func TestFetchColumnsMetadataScanRowError(t *testing.T) {
	db, mock, err := sqlmock.New()

	require.Nil(t, err)
	defer db.Close()

	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"}, db)

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

	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"}, db)
	sqlErr := errors.New("rows error")

	rows := sqlmock.NewRows([]string{
		"COLUMN_NAME", "DATA_TYPE", "NULLABLE", "DATA_LENGTH", "DATA_PRECISION", "DATA_SCALE",
	}).AddRow("TOTAL", "NUMBER", "F", 22, 17, 2).RowError(0, sqlErr)

	mock.ExpectQuery("^SELECT (.+) FROM  ALL_TAB_COLS").WillReturnRows(rows)

	_, err = r.FetchColumnsMetadata(schema.Table{Name: "customers", Ignore: []schema.Ignore{"id"}})
	assert.ErrorIs(t, err, sqlErr)
}

func TestFetchData(t *testing.T) {
	r, _ := reader.NewReader(reader.DataSource{DBMSName: "oracle"}, nil)
	assert.Panics(
		t,
		func() { r.FetchData("", []reader.DBColumn{}, []string{}, [][]interface{}{}) },
		"not implemented yet",
	)
}
