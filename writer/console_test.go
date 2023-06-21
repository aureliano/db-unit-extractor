package writer_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/writer"
	"github.com/stretchr/testify/assert"
)

func TestConsoleWriteHeader(t *testing.T) {
	w := writer.ConsoleWriter{}
	assert.Nil(t, w.WriteHeader())
}

func TestConsoleWriteFooter(t *testing.T) {
	w := writer.ConsoleWriter{}
	assert.Nil(t, w.WriteFooter())
}

func TestConsoleWrite(t *testing.T) {
	w := writer.ConsoleWriter{}

	rows := make([][]*reader.DBColumn, 2)
	rows[0] = []*reader.DBColumn{{Name: "id", Value: 1}, {}}
	rows[1] = []*reader.DBColumn{{}, {Name: "name", Value: "Giovanni Pierluigi da Palestrina"}}

	err := w.Write("", rows)
	assert.Nil(t, err)
}
