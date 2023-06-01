package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type XMLWriter struct {
	Formatted bool
	Directory string
	Name      string
	file      *os.File
}

func (w *XMLWriter) WriteHeader() error {
	err := os.MkdirAll(w.Directory, os.ModePerm)
	if err != nil {
		return err
	}

	path := filepath.Join(w.Directory, fmt.Sprintf("%s.xml", w.Name))
	w.file, err = os.Create(path)
	if err != nil {
		return err
	}

	_, err = w.file.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><dataset>"))
	if err != nil {
		_ = w.file.Close()
	}

	return err
}

func (w *XMLWriter) WriteFooter() error {
	_, err := w.file.Write([]byte("</dataset>"))
	if err != nil {
		_ = w.file.Close()
		return err
	}

	return w.file.Close()
}

func (w *XMLWriter) Write(table string, rows []map[string]interface{}) error {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("<%s", table))

	for _, row := range rows {
		for name, value := range row {
			sb.WriteString(fmt.Sprintf(" %s=\"%v\"", name, value))
		}
	}

	sb.WriteString("/>")
	_, err := w.file.Write([]byte(sb.String()))

	if err != nil {
		_ = w.file.Close()
	}

	return err
}
