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

	_, err = w.file.Write(fileHeader(w.Formatted))
	if err != nil {
		_ = w.file.Close()
	}

	return err
}

func (w *XMLWriter) WriteFooter() error {
	_, err := w.file.Write(fileFooter())
	if err != nil {
		_ = w.file.Close()
		return err
	}

	return w.file.Close()
}

func (w *XMLWriter) Write(table string, rows []map[string]interface{}) error {
	_, err := w.file.Write(fileBody(w.Formatted, table, rows))
	if err != nil {
		_ = w.file.Close()
	}

	return err
}

func fileHeader(formatted bool) []byte {
	if formatted {
		return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<dataset>\n")
	}

	return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><dataset>")
}

func fileFooter() []byte {
	return []byte("</dataset>")
}

func fileBody(formatted bool, table string, rows []map[string]interface{}) []byte {
	if formatted {
		return formattedXMLRecord(table, rows)
	}

	return unformattedXMLRecord(table, rows)
}

func formattedXMLRecord(table string, rows []map[string]interface{}) []byte {
	sb := strings.Builder{}

	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("  <%s", table))

		for name, value := range row {
			if value != nil {
				sb.WriteString(fmt.Sprintf("\n    %s=\"%v\"", name, value))
			}
		}
		sb.WriteString("/>\n")
	}

	return []byte(sb.String())
}

func unformattedXMLRecord(table string, rows []map[string]interface{}) []byte {
	sb := strings.Builder{}

	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("<%s", table))
		for name, value := range row {
			if value != nil {
				sb.WriteString(fmt.Sprintf(" %s=\"%v\"", name, value))
			}
		}
		sb.WriteString("/>")
	}

	return []byte(sb.String())
}
