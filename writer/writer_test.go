package writer_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/writer"
	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	_, err := writer.NewWriter(writer.FileConf{})
	assert.ErrorIs(t, err, writer.ErrUnsupportedFileWriter)
}

func TestNewWriterConsole(t *testing.T) {
	w, err := writer.NewWriter(writer.FileConf{Type: "console"})
	assert.Nil(t, err)
	assert.IsType(t, &writer.ConsoleWriter{}, w)
}

func TestNewWriterXML(t *testing.T) {
	w, err := writer.NewWriter(writer.FileConf{Type: "xml"})
	assert.Nil(t, err)
	assert.IsType(t, &writer.XMLWriter{}, w)
}

func TestSupportedTypes(t *testing.T) {
	types := writer.SupportedTypes()
	assert.Len(t, types, 2)
	assert.Equal(t, "console", types[0])
	assert.Equal(t, "xml", types[1])
}
