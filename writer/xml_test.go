package writer_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/writer"
	"github.com/stretchr/testify/assert"

	"github.com/agiledragon/gomonkey/v2"
)

func TestXMLWriteHeaderMkdirAllError(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.MkdirAll, func(string, fs.FileMode) error {
		return fmt.Errorf("mkdir error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{}

	assert.Equal(t, "mkdir error", w.WriteHeader().Error())
}

func TestXMLWriteHeaderFileCreationError(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.OpenFile, func(string, int, fs.FileMode) (*os.File, error) {
		return nil, fmt.Errorf("file creation error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file creation error", w.WriteHeader().Error())
}

func TestXMLWriteHeaderFileWritingError(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&os.File{}, "Write", func([]byte) (int, error) {
		return 0, fmt.Errorf("file writing error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file writing error", w.WriteHeader().Error())
}

func TestXMLWriteFooterFileWritingError(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&os.File{}, "Write", func([]byte) (int, error) {
		return 0, fmt.Errorf("file writing error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file writing error", w.WriteFooter().Error())
}

func TestXMLWriteBodyFileWritingError(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&os.File{}, "Write", func([]byte) (int, error) {
		return 0, fmt.Errorf("file writing error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file writing error", w.Write("", make([][]*reader.DBColumn, 2)).Error())
}

func TestXMLWriteUnformattedEmptyData(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.XMLWriter{
		Formatted: false,
		Directory: dir,
		Name:      "test-write-unformatted",
	}

	assert.Nil(t, w.WriteHeader())

	assert.Nil(t, w.Write("products", [][]*reader.DBColumn{}))

	assert.Nil(t, w.WriteFooter())

	bytes, _ := os.ReadFile(filepath.Join(dir, fmt.Sprintf("%s.xml", w.Name)))
	assert.Equal(t, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><dataset></dataset>", string(bytes))
}

func TestXMLWriteUnformatted(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.XMLWriter{
		Formatted: false,
		Directory: dir,
		Name:      "test-write-unformatted",
	}

	assert.Nil(t, w.WriteHeader())

	rows := make([][]*reader.DBColumn, 2)
	rows[0] = []*reader.DBColumn{
		{Name: "id", Value: 1},
		{Name: "name", Value: "shirt"},
		{Name: "description", Value: "black shirt"},
		{Name: "price", Value: 14.50},
	}

	assert.Nil(t, w.Write("products", rows))

	assert.Nil(t, w.WriteFooter())

	bytes, _ := os.ReadFile(filepath.Join(dir, fmt.Sprintf("%s.xml", w.Name)))
	xml := string(bytes)

	assert.Contains(t, xml, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><dataset><products ")
	assert.Contains(t, xml, "id=\"1\"")
	assert.Contains(t, xml, "name=\"shirt\"")
	assert.Contains(t, xml, "description=\"black shirt\"")
	assert.Contains(t, xml, "price=\"14.5\"")
	assert.Contains(t, xml, "/></dataset>")
}

func TestXMLWriteFormatted(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.XMLWriter{
		Formatted: true,
		Directory: dir,
		Name:      "test-write-formatted",
	}

	assert.Nil(t, w.WriteHeader())

	rows := make([][]*reader.DBColumn, 2)
	rows[0] = []*reader.DBColumn{
		{Name: "id", Value: 1},
		{Name: "name", Value: "shirt"},
		{Name: "description", Value: "black shirt"},
		{Name: "price", Value: 14.50},
	}

	assert.Nil(t, w.Write("products", rows))

	assert.Nil(t, w.WriteFooter())

	bytes, _ := os.ReadFile(filepath.Join(dir, fmt.Sprintf("%s.xml", w.Name)))
	xml := string(bytes)

	assert.Contains(t, xml, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<dataset>\n  <products\n")
	assert.Contains(t, xml, "    id=\"1\"")
	assert.Contains(t, xml, "    name=\"shirt\"")
	assert.Contains(t, xml, "    description=\"black shirt\"")
	assert.Contains(t, xml, "    price=\"14.5\"")
	assert.Contains(t, xml, "/>\n</dataset>")
}
