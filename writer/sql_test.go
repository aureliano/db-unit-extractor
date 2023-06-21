package writer_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
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

func TestWriteHeader(t *testing.T) {
	//w := writer.SQLWriter{}
	//assert.Nil(t, w.WriteHeader())
}

func TestWriteFooter(t *testing.T) {
	w := writer.SQLWriter{}
	assert.Nil(t, w.WriteFooter())
}

func TestWrite(t *testing.T) {
	w := writer.SQLWriter{}
	assert.Nil(t, w.Write("", nil))
}
