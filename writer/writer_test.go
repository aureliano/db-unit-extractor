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

func TestNewWriterDummy(_ *testing.T) {
	_, _ = writer.NewWriter(writer.FileConf{Type: "dummy"})
}

func TestWriteHeader(_ *testing.T) {
	w := writer.DummyWriter{}
	w.WriteHeader()
}

func TestWriteFooter(_ *testing.T) {
	w := writer.DummyWriter{}
	w.WriteFooter()
}

func TestWrite(_ *testing.T) {
	w := writer.DummyWriter{}
	_ = w.Write("", nil)
}
