package dataconv_test

import (
	"testing"
	"time"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/stretchr/testify/assert"
)

func TestConvertTimeSourceIsNotTime(t *testing.T) {
	c := dataconv.DateTimeISO8601Converter{}
	source := "2023-06-09T14:31:16.478 -0300"
	actual, err := c.Convert(source)

	assert.Equal(t, "'2023-06-09T14:31:16.478 -0300' is not time.Time", err.Error())
	assert.Nil(t, actual)
}

func TestConvertDateAndTime(t *testing.T) {
	c := dataconv.DateTimeISO8601Converter{}
	location, _ := time.LoadLocation("America/Sao_Paulo")
	source := time.Date(2023, time.June, 9, 14, 31, 16, 478587456, location)
	expected := "2023-06-09T14:31:16.478 -0300"
	actual, err := c.Convert(source)

	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestConvertDateOnly(t *testing.T) {
	c := dataconv.DateTimeISO8601Converter{}
	location, _ := time.LoadLocation("America/Sao_Paulo")
	source := time.Date(2023, time.June, 19, 0, 0, 0, 0, location)
	expected := "2023-06-19"
	actual, err := c.Convert(source)

	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestHandleTime(t *testing.T) {
	c := dataconv.DateTimeISO8601Converter{}
	assert.False(t, c.Handle("2023-06-09T14:31:16.478 -0300"))
	assert.True(t, c.Handle(time.Now()))
}
