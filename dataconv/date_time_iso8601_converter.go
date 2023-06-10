package dataconv

import (
	"fmt"
	"time"
)

type DateTimeISO8601Converter struct{}

const dateTimeISO8601Layout = "2006-01-02T15:04:05.999 -0700"

func (DateTimeISO8601Converter) Convert(source interface{}) (interface{}, error) {
	tm, isTime := source.(time.Time)
	if !isTime {
		return nil, fmt.Errorf("'%v' is not time.Time", source)
	}

	return tm.Format(dateTimeISO8601Layout), nil
}
