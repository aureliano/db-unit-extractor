package writer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type SQLWriter struct {
	Formatted bool
	Directory string
	Name      string
	file      *os.File
}

func (w *SQLWriter) WriteHeader() error {
	err := os.MkdirAll(w.Directory, os.ModePerm)
	if err != nil {
		log.Printf("SQL.WriteHeader\nMake directory %s failed with `%s'\n", w.Directory, err.Error())
		return err
	}

	path := filepath.Join(w.Directory, fmt.Sprintf("%s.sql", w.Name))
	w.file, err = os.Create(path)
	if err != nil {
		log.Printf("SQL.WriteHeader\nFile %s not created: `%s'\n", path, err.Error())
		return err
	}

	return nil
}

func (w *SQLWriter) WriteFooter() error {
	return w.file.Close()
}

func (w *SQLWriter) Write(table string, rows []map[string]interface{}) error { return nil }
