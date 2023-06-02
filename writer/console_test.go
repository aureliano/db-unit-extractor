package writer_test

import (
	"testing"

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

	rows := make([]map[string]interface{}, 2)

	row := make(map[string]interface{})
	row["id"] = 1
	rows[0] = row

	row = make(map[string]interface{})
	row["name"] = "Giovanni Pierluigi da Palestrina"
	rows[1] = row

	err := w.Write("", rows)
	assert.Nil(t, err)
}
