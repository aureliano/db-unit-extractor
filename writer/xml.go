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
	_, err := w.file.Write(fileFooter(w.Formatted))
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

func fileFooter(formatted bool) []byte {
	if formatted {
		return []byte("\n</dataset>\n")
	}

	return []byte("</dataset>")
}

func fileBody(formatted bool, table string, rows []map[string]interface{}) []byte {
	sb := strings.Builder{}

	if formatted {
		sb.WriteString(fmt.Sprintf("  <%s\n", table))

		for _, row := range rows {
			li := len(row) - 1
			i := 0
			for name, value := range row {
				sb.WriteString(fmt.Sprintf("    %s=\"%v\"", name, value))
				if i < li {
					sb.WriteString("\n")
				}
				i++
			}
		}
	} else {
		sb.WriteString(fmt.Sprintf("<%s", table))

		for _, row := range rows {
			for name, value := range row {
				sb.WriteString(fmt.Sprintf(" %s=\"%v\"", name, value))
			}
		}
	}

	sb.WriteString("/>")

	return []byte(sb.String())
}
