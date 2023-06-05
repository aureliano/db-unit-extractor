package main

import (
	"os"
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

	_, err := parseYAML(args[0])
	if err != nil {
		panic(err)
	}

	_, err = parseXML(args[1])
	if err != nil {
		panic(err)
	}
}
