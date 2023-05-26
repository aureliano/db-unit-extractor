package writer

import (
	"fmt"
	"os"
)

type ConsoleWriter struct{}

func (ConsoleWriter) WriteHeader() {
	_ = 0
}

func (ConsoleWriter) WriteFooter() {
	_ = 0
}

func (ConsoleWriter) Write(table string, rows []map[string]interface{}) error {
	fmt.Fprintln(os.Stdout, " >", table)
	for _, row := range rows {
		for name, value := range row {
			fmt.Fprintf(os.Stdout, "   %s: %v\n", name, value)
		}
	}

	return nil
}
