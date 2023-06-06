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

		allMatch, msg := fieldsMatch(t, actual, ind)
		if !allMatch {
			return fmt.Errorf("table %s doesn't match: %v != %v\n%s", t.Name, t, actual.Tables[ind[0]], msg)
		}
	}

	return nil
}

func fieldsMatch(t Table, actual *DataSet, ind []int) (bool, string) {
	msg := strings.Builder{}
	var allMatch bool

	for _, i := range ind {
		allMatch = true
		if len(t.Fields) != len(actual.Tables[i].Fields) {
			continue
		}

		for _, f := range t.Fields {
			match := fieldMatch(f, actual.Tables[i].Fields)
			if !match {
				msg.WriteString(fmt.Sprintf("'%s' => '%s'\n", f.Name, f.Value))
			}

			allMatch = allMatch && match
		}

		if allMatch {
			break
		}
	}

	return allMatch, msg.String()
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
