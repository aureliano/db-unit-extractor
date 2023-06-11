package dataconv_test

import (
	"encoding/base64"
	"testing"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/stretchr/testify/assert"
)

func TestConvertBytesSourceIsNotBytes(t *testing.T) {
	c := dataconv.BlobConverter{}
	source := "test"
	actual, err := c.Convert(source)

	assert.Equal(t, "'test' is not []byte", err.Error())
	assert.Nil(t, actual)
}

func TestConvertBytes(t *testing.T) {
	c := dataconv.BlobConverter{}
	source := []byte("test bytes converter")
	expected := "dGVzdCBieXRlcyBjb252ZXJ0ZXI="
	actual, err := c.Convert(source)

	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)

	bts, _ := base64.StdEncoding.DecodeString(actual.(string))
	assert.Equal(t, source, bts)
}

func TestHandleBlob(t *testing.T) {
	c := dataconv.BlobConverter{}
	assert.False(t, c.Handle(456))
	assert.True(t, c.Handle([]byte{1, 3, 65}))
}
