package writer

import (
	"fmt"
	"os"
)

type ConsoleWriter struct{}

func (*ConsoleWriter) WriteHeader() error {
	return nil
}

func (*ConsoleWriter) WriteFooter() error {
	return nil
}

func (*ConsoleWriter) Write(table string, rows []map[string]interface{}) error {
	for _, row := range rows {
		fmt.Fprintln(os.Stdout, " >", table)

		for name, value := range row {
			if value != nil {
				fmt.Fprintf(os.Stdout, "   %s: %v\n", name, value)
			}
		}
	}

	fmt.Fprintln(os.Stdout)

	return nil
}
