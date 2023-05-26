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
