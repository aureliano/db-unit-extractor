package main

import (
	"fmt"
	"os"
	"strings"
)

type Field struct {
	Name  string
	Value string
}

type Table struct {
	Name   string
	Fields []Field
}

type DataSet struct {
	Tables []Table
}

func main() {
	const numExpectedParams = 2
	args := os.Args[1:]

	if len(args) != numExpectedParams {
		panic("Expected two arguments: dataset and expectation.")
	}

	expected, err := parseYAML(args[0])
	if err != nil {
		panic(err)
	}

	actual, err := parseXML(args[1])
	if err != nil {
		panic(err)
	}

	if err = compare(expected, actual); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func compare(expected, actual *DataSet) error {
	l1 := len(expected.Tables)
	l2 := len(actual.Tables)
	if l1 != l2 {
		return fmt.Errorf("len(expected.Tables) != len(actual.Tables) %d != %d", l1, l2)
	}

	for _, t := range expected.Tables {
		ind := indexOfTable(t, actual)
		if len(ind) == 0 {
			return fmt.Errorf("table %s not found", t.Name)
		}

		anyMatch := false
		for _, i := range ind {
			if len(t.Fields) != len(actual.Tables[i].Fields) {
				continue
			}

			for _, f := range t.Fields {
				if fieldMatch(f, actual.Tables[i].Fields) {
					anyMatch = true
					break
				}
			}
		}

		if !anyMatch {
			return fmt.Errorf("table %s doesn't match: %v != %v", t.Name, t, actual.Tables[ind[0]])
		}
	}

	return nil
}

func indexOfTable(t Table, ds *DataSet) []int {
	ind := make([]int, 0)
	for i, table := range ds.Tables {
		if t.Name == table.Name {
			ind = append(ind, i)
		}
	}

	return ind
}

func fieldMatch(f Field, fields []Field) bool {
	for _, field := range fields {
		if strings.EqualFold(f.Name, field.Name) && f.Value == field.Value {
			return true
		}
	}

	return false
}
