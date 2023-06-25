package main

import (
	"os"
	"regexp"
	"strings"
)

var (
	insertRegExp = regexp.MustCompile(`(insert\sinto|INSERT\sINTO)\s(\w+)\(([^)]+)\)\s+(values|VALUES)\s*([^;]+)`)
	spacesRegExp = regexp.MustCompile(`\s+`)
)

func parseSQL(path string) (*DataSet, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	sql := string(byteValue)
	inserts := strings.Split(sql, ";")
	tables := make([]Table, 0, len(inserts))

	for _, insert := range inserts {
		res := insertRegExp.FindAllStringSubmatch(insert, -1)
		if len(res) == 0 {
			continue
		}
		tname := res[0][2]
		fnames := strings.Split(res[0][3], ",")
		arrvalues := findInsertValues(res[0][5])

		fields := make([]Field, len(fnames))
		for j, fname := range fnames {
			fields[j] = Field{Name: spacesRegExp.ReplaceAllString(fname, "")}
		}

		for _, values := range arrvalues {
			cpFields := make([]Field, len(fnames))
			copy(cpFields, fields)

			for i, value := range values {
				field := &cpFields[i]
				field.Value = value
			}

			table := Table{Name: tname, Fields: cpFields}
			tables = append(tables, table)
		}
	}

	return &DataSet{Tables: tables}, nil
}

func findInsertValues(text string) [][]string {
	re := regexp.MustCompile(`\)\s*\(`)
	fre := regexp.MustCompile(`'([^']+)'`)
	indexes := re.FindAllStringIndex(text, -1)
	values := make([][]string, len(indexes)+1)
	lines := make([]string, len(indexes)+1)

	if len(indexes) == 0 {
		lines[0] = text
	} else {
		lines[0] = text[:indexes[0][0]]

		size := len(indexes)
		for i := 0; i < size; i++ {
			if i+1 < size {
				begin := indexes[i][0]
				end := indexes[i+1][0]
				lines[i+1] = text[begin:end]
			}
		}

		index := indexes[size-1][0]
		lines[len(lines)-1] = text[index:]
	}

	for i, line := range lines {
		fields := fre.FindAllStringSubmatch(line, -1)
		source := make([]string, len(fields))

		for j, value := range fields {
			source[j] = value[1]
		}

		values[i] = source
	}

	return values
}
