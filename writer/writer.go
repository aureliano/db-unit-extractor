package writer

import (
	"errors"
	"fmt"
	"strings"
)

type FileConf struct {
	Type      string
	Formatted bool
	Directory string
}

var ErrUnsupportedFileWriter = errors.New("unsupported file type")

type FileWriter interface {
	WriteHeader()
	WriteFooter()
	Write(table string, rows []map[string]interface{}) error
}

func NewWriter(conf FileConf) (FileWriter, error) {
	switch {
	case strings.EqualFold(conf.Type, "console"):
		return ConsoleWriter{}, nil
	case strings.EqualFold(conf.Type, "xml"):
		return XMLWriter{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileWriter, conf.Type)
	}
}

func SupportedTypes() []string {
	return []string{"console", "xml"}
}
