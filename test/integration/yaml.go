package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type YAMLDataSet struct {
	Tables []YAMLTable `yaml:"tables"`
}

type YAMLTable struct {
	Name    string `yaml:"name"`
	Records []Tag
}

func parseYAML(path string) (*DataSet, error) {
	schema := make(map[interface{}]interface{})
	yml, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	if err = yaml.UnmarshalStrict(yml, &schema); err != nil {
		return nil, err
	}

	tables, _ := schema["tables"].([]interface{})

	dataset := DataSet{
		Tables: make([]Table, len(tables)),
	}

	for i, table := range tables {
		t, _ := table.(map[interface{}]interface{})
		name := t["name"]
		records, _ := t["records"].([]interface{})
		fields := make([]Field, len(records))

		for j, rec := range records {
			var fname, fvalue string
			for k, v := range rec.(map[interface{}]interface{}) {
				fname, _ = k.(string)
				fvalue = fmt.Sprintf("%v", v)
			}

			fields[j] = Field{
				Name:  fname,
				Value: fvalue,
			}
		}

		dataset.Tables[i] = Table{
			Name:   name.(string),
			Fields: fields,
		}
	}

	return &dataset, nil
}
