package writer_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/writer"
	"github.com/stretchr/testify/assert"
)

func TestSQLWriteHeaderMkdirAllError(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.MkdirAll, func(string, fs.FileMode) error {
		return fmt.Errorf("mkdir error")
	})
	defer patches.Reset()

	w := writer.SQLWriter{}

	assert.Equal(t, "mkdir error", w.WriteHeader().Error())
}

func TestSQLWriteHeaderFileCreationError(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.OpenFile, func(string, int, fs.FileMode) (*os.File, error) {
		return nil, fmt.Errorf("file creation error")
	})
	defer patches.Reset()

	w := writer.SQLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file creation error", w.WriteHeader().Error())
}

func TestSQLWriteBodyFileWritingError(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&os.File{}, "Write", func([]byte) (int, error) {
		return 0, fmt.Errorf("file writing error")
	})
	defer patches.Reset()

	w := writer.SQLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file writing error", w.Write("", [][]*reader.DBColumn{{}}).Error())
}

func TestSQLWriteUnformattedEmptyData(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.SQLWriter{
		Formatted: false,
		Directory: dir,
		Name:      "test-write-unformatted",
	}

	assert.Nil(t, w.WriteHeader())

	assert.Nil(t, w.Write("products", [][]*reader.DBColumn{}))

	assert.Nil(t, w.WriteFooter())

	bytes, _ := os.ReadFile(filepath.Join(dir, fmt.Sprintf("%s.sql", w.Name)))
	assert.Empty(t, string(bytes))
}

func TestSQLWriteUnformatted(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.SQLWriter{
		Formatted: false,
		Directory: dir,
		Name:      "test-write-unformatted",
	}

	assert.Nil(t, w.WriteHeader())

	rows := [][]*reader.DBColumn{{
		{Name: "id", Value: 1},
		{Name: "name", Value: "shirt"},
		{Name: "description", Value: "black shirt"},
		{Name: "price", Value: 14.50},
	}, {
		{Name: "id", Value: 2},
		{Name: "name", Value: "pant"},
		{Name: "description"},
		{Name: "price", Value: 26.35},
	}}

	assert.Nil(t, w.Write("products", rows))

	assert.Nil(t, w.WriteFooter())

	bytes, _ := os.ReadFile(filepath.Join(dir, fmt.Sprintf("%s.sql", w.Name)))
	actual := string(bytes)
	expected := "insert into products(id,name,description,price) " +
		"values('1','shirt','black shirt','14.5')('2','pant',null,'26.35');"

	assert.Equal(t, expected, actual)
}

func TestSQLWriteFormatted(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.SQLWriter{
		Formatted: true,
		Directory: dir,
		Name:      "test-write-unformatted",
	}

	assert.Nil(t, w.WriteHeader())

	rows := [][]*reader.DBColumn{{
		{Name: "id", Value: 1},
		{Name: "name", Value: "shirt"},
		{Name: "description", Value: "black shirt"},
		{Name: "price", Value: 14.50},
	}, {
		{Name: "id", Value: 2},
		{Name: "name", Value: "pant"},
		{Name: "description"},
		{Name: "price", Value: 26.35},
	}}

	assert.Nil(t, w.Write("products", rows))

	assert.Nil(t, w.WriteFooter())

	bytes, _ := os.ReadFile(filepath.Join(dir, fmt.Sprintf("%s.sql", w.Name)))
	actual := string(bytes)
	expected := "INSERT INTO products(\n  id, name, description, price\n) " +
		"VALUES\n  ('1', 'shirt', 'black shirt', '14.5')\n  ('2', 'pant', null, '26.35');"

	assert.Equal(t, expected, actual)
}
