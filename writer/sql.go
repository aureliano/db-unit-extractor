package writer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aureliano/db-unit-extractor/reader"
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

func (w *SQLWriter) Write(table string, rows [][]*reader.DBColumn) error {
	if len(rows) == 0 {
		return nil
	}

	content := sqlFileBody(w.Formatted, table, rows)
	_, err := w.file.Write(content)
	if err != nil {
		log.Printf("SQL.Write\nWriting to file failed: `%s'\nContent: %s\n", err.Error(), content)
		_ = w.file.Close()
	}

	return err
}

func sqlFileBody(formatted bool, table string, rows [][]*reader.DBColumn) []byte {
	if formatted {
		return formattedSQLRecord(table, rows)
	}

	return unformattedSQLRecord(table, rows)
}

func formattedSQLRecord(table string, rows [][]*reader.DBColumn) []byte {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("INSERT INTO %s(\n  ", table))

	for i, cname := range rows[0] {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(cname.Name)
	}

	sb.WriteString("\n) VALUES")
	for _, row := range rows {
		sb.WriteString("\n  (")
		for i, column := range row {
			if i > 0 {
				sb.WriteString(", ")
			}

			if column.Value == nil {
				sb.WriteString("null")
			} else {
				sb.WriteString(fmt.Sprintf("'%v'", column.Value))
			}
		}
		sb.WriteString(")")
	}
	sb.WriteString(";\n")

	return []byte(sb.String())
}

func unformattedSQLRecord(table string, rows [][]*reader.DBColumn) []byte {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("insert into %s(", table))

	for i, cname := range rows[0] {
		if i > 0 {
			sb.WriteRune(',')
		}

		sb.WriteString(cname.Name)
	}

	sb.WriteString(") values")
	for _, row := range rows {
		sb.WriteRune('(')
		for i, column := range row {
			if i > 0 {
				sb.WriteRune(',')
			}

			if column.Value == nil {
				sb.WriteString("null")
			} else {
				sb.WriteString(fmt.Sprintf("'%v'", column.Value))
			}
		}
		sb.WriteRune(')')
	}
	sb.WriteRune(';')

	return []byte(sb.String())
}
