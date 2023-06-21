package writer_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

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
	assert.Equal(t, "file writing error", w.Write("", make([]map[string]interface{}, 2)).Error())
}

func TestXMLWriteUnformatted(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.XMLWriter{
		Formatted: false,
		Directory: dir,
		Name:      "test-write-unformatted",
	}

	assert.Nil(t, w.WriteHeader())

	row := make(map[string]interface{})
	row["id"] = 1
	row["name"] = "shirt"
	row["description"] = "black shirt"
	row["price"] = 14.50

	rows := make([]map[string]interface{}, 2)
	rows[0] = row

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

	row := make(map[string]interface{})
	row["id"] = 1
	row["name"] = "shirt"
	row["description"] = "black shirt"
	row["price"] = 14.50

	rows := make([]map[string]interface{}, 2)
	rows[0] = row

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
