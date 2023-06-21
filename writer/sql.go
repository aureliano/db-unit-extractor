package writer

import "os"

type SQLWriter struct {
	Formatted bool
	Directory string
	Name      string
	file      *os.File
}

func (w *SQLWriter) WriteHeader() error { return nil }

func (w *SQLWriter) WriteFooter() error { return nil }

func (w *SQLWriter) Write(table string, rows []map[string]interface{}) error { return nil }
