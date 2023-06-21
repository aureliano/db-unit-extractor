package dataconv

import (
	"fmt"
	"time"
)

type DateTimeISO8601Converter struct{}

const (
	dateISO8601Layout     = "2006-01-02"
	dateTimeISO8601Layout = "2006-01-02T15:04:05.999 -0700"
)

func (DateTimeISO8601Converter) Convert(_ string, source interface{}) (interface{}, error) {
	tm, isTime := source.(time.Time)
	if !isTime {
		return nil, fmt.Errorf("'%v' is not time.Time", source)
	}

	if isDateAndTime(tm) {
		return tm.Format(dateTimeISO8601Layout), nil
	}

	return tm.Format(dateISO8601Layout), nil
}

func (DateTimeISO8601Converter) Handle(vl interface{}) bool {
	_, handled := vl.(time.Time)
	return handled
}

func isDateAndTime(tm time.Time) bool {
	return !(tm.Hour() == 0 && tm.Minute() == 0 && tm.Second() == 0)
}
