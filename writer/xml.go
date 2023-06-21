package writer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aureliano/db-unit-extractor/reader"
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
		log.Printf("XML.WriteHeader\nMake directory %s failed with `%s'\n", w.Directory, err.Error())
		return err
	}

	path := filepath.Join(w.Directory, fmt.Sprintf("%s.xml", w.Name))
	w.file, err = os.Create(path)
	if err != nil {
		log.Printf("XML.WriteHeader\nFile %s not created: `%s'\n", path, err.Error())
		return err
	}

	content := xmlFileHeader(w.Formatted)
	_, err = w.file.Write(content)
	if err != nil {
		log.Printf("XML.WriteHeader\nWriting to file failed: `%s'\nContent: %s\n", err.Error(), content)
		_ = w.file.Close()
	}

	return err
}

func (w *XMLWriter) WriteFooter() error {
	content := xmlFileFooter()
	_, err := w.file.Write(content)
	if err != nil {
		log.Printf("XML.WriteFooter\nWriting to file failed: `%s'\nContent: %s\n", err.Error(), content)
		_ = w.file.Close()
		return err
	}

	return w.file.Close()
}

func (w *XMLWriter) Write(table string, rows [][]*reader.DBColumn) error {
	if len(rows) == 0 {
		return nil
	}

	content := xmlFileBody(w.Formatted, table, rows)
	_, err := w.file.Write(content)
	if err != nil {
		log.Printf("XML.Write\nWriting to file failed: `%s'\nContent: %s\n", err.Error(), content)
		_ = w.file.Close()
	}

	return err
}

func xmlFileHeader(formatted bool) []byte {
	if formatted {
		return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<dataset>\n")
	}

	return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><dataset>")
}

func xmlFileFooter() []byte {
	return []byte("</dataset>")
}

func xmlFileBody(formatted bool, table string, rows [][]*reader.DBColumn) []byte {
	if formatted {
		return formattedXMLRecord(table, rows)
	}

	return unformattedXMLRecord(table, rows)
}

func formattedXMLRecord(table string, rows [][]*reader.DBColumn) []byte {
	sb := strings.Builder{}

	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("  <%s", table))

		for _, column := range row {
			if column.Value != nil {
				sb.WriteString(fmt.Sprintf("\n    %s=\"%v\"", column.Name, column.Value))
			}
		}
		sb.WriteString("/>\n")
	}

	return []byte(sb.String())
}

func unformattedXMLRecord(table string, rows [][]*reader.DBColumn) []byte {
	sb := strings.Builder{}

	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("<%s", table))
		for _, column := range row {
			if column.Value != nil {
				sb.WriteString(fmt.Sprintf(" %s=\"%v\"", column.Name, column.Value))
			}
		}
		sb.WriteString("/>")
	}

	return []byte(sb.String())
}
