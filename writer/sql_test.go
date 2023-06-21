package writer_test

import (
	"testing"

	"github.com/aureliano/db-unit-extractor/writer"
	"github.com/stretchr/testify/assert"
)

func TestWriteHeader(t *testing.T) {
	w := writer.SQLWriter{}
	assert.Nil(t, w.WriteHeader())
}

func TestWriteFooter(t *testing.T) {
	w := writer.SQLWriter{}
	assert.Nil(t, w.WriteFooter())
}

func TestWrite(t *testing.T) {
	w := writer.SQLWriter{}
	assert.Nil(t, w.Write("", nil))
}
