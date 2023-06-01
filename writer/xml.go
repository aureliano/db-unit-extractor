package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

type XMLWriter struct {
	Formatted bool
	Directory string
	Name      string
	file      *os.File
}

func (w XMLWriter) WriteHeader() error {
	err := os.MkdirAll(w.Directory, os.ModePerm)
	if err != nil {
		return err
	}

	path := filepath.Join(w.Directory, fmt.Sprintf("%s.xml", w.Name))
	w.file, err = os.Create(path)
	if err != nil {
		return err
	}

	return nil
}

func (XMLWriter) WriteFooter() error {
	return nil
}

func (XMLWriter) Write(table string, rows []map[string]interface{}) error {
	return nil
}
