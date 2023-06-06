package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

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
		Tables: make([]Table, 0, len(tables)),
	}

	for _, table := range tables {
		t, _ := table.(map[interface{}]interface{})
		name := t["name"]
		records, _ := t["records"].([]interface{})

		for _, record := range records {
			r, _ := record.(map[interface{}]interface{})
			fields := make([]Field, len(r))

			var fname, fvalue string
			l := 0
			for k, v := range r {
				fname, _ = k.(string)
				fvalue = fmt.Sprintf("%v", v)

				fields[l] = Field{
					Name:  fname,
					Value: fvalue,
				}
				l++
			}

			dataset.Tables = append(dataset.Tables, Table{
				Name:   name.(string),
				Fields: fields,
			})
		}
	}

	return &dataset, nil
}
