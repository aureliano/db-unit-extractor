package writer

import (
	"fmt"
	"os"

	"github.com/aureliano/db-unit-extractor/reader"
)

type ConsoleWriter struct{}

func (*ConsoleWriter) WriteHeader() error {
	return nil
}

func (*ConsoleWriter) WriteFooter() error {
	return nil
}

func (*ConsoleWriter) Write(table string, rows [][]*reader.DBColumn) error {
	for _, row := range rows {
		fmt.Fprintln(os.Stdout, " >", table)

		for _, column := range row {
			if column.Value != nil {
				fmt.Fprintf(os.Stdout, "   %s: %v\n", column.Name, column.Value)
			}
		}
	}

	fmt.Fprintln(os.Stdout)

	return nil
}
