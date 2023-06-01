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

func TestWriteHeaderMkdirAllError(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.MkdirAll, func(string, fs.FileMode) error {
		return fmt.Errorf("mkdir error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{}

	assert.Equal(t, "mkdir error", w.WriteHeader().Error())
}

func TestWriteHeaderFileCreationError(t *testing.T) {
	patches := gomonkey.ApplyFunc(os.OpenFile, func(string, int, fs.FileMode) (*os.File, error) {
		return nil, fmt.Errorf("file creation error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file creation error", w.WriteHeader().Error())
}

func TestWriteHeaderFileWritingError(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&os.File{}, "Write", func([]byte) (int, error) {
		return 0, fmt.Errorf("file writing error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file writing error", w.WriteHeader().Error())
}

func TestWriteFooterFileWritingError(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&os.File{}, "Write", func([]byte) (int, error) {
		return 0, fmt.Errorf("file writing error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file writing error", w.WriteFooter().Error())
}

func TestWriteBodyFileWritingError(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&os.File{}, "Write", func([]byte) (int, error) {
		return 0, fmt.Errorf("file writing error")
	})
	defer patches.Reset()

	w := writer.XMLWriter{Directory: filepath.Join(os.TempDir(), "db-unit-extractor", "writer")}
	assert.Equal(t, "file writing error", w.Write("", make([]map[string]interface{}, 2)).Error())
}

func TestWriteUnformatted(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "db-unit-extractor", "writer")
	w := writer.XMLWriter{
		Formatted: false,
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
}
