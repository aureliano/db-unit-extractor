package writer

import (
	"errors"
	"fmt"
)

type FileConf struct {
	Type      string
	Formatted bool
}

var ErrUnsupportedFileWriter = errors.New("unsupported file type")

type FileWriter interface {
	WriteHeader()
	WriteFooter()
	Write(table string, rows []map[string]interface{}) error
}

type DummyWriter struct{}

func (DummyWriter) WriteHeader() { _ = 0 }

func (DummyWriter) WriteFooter() { _ = 0 }

func (DummyWriter) Write(_ string, _ []map[string]interface{}) error {
	return nil
}

func NewWriter(conf FileConf) (FileWriter, error) {
	if conf.Type == "dummy" {
		return DummyWriter{}, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileWriter, conf.Type)
}
