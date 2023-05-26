package writer

import (
	"errors"
	"fmt"
)

type FileConf struct {
	Type                 string
	Formatted            bool
	TabSize              int
	BreakLineAfterColumn bool
}

var ErrUnsupportedFileWriter = errors.New("unsupported file type")

type FileWriter interface {
	WriteHeader()
	WriteFooter()
	Write(table string, rows []map[string]interface{}) error
}

func NewWriter(conf FileConf) (FileWriter, error) {
	return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileWriter, conf.Type)
}
