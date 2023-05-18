package extractor

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Schema struct {
	Converters []string `yaml:"converters"`
	Tables     []struct {
		Name    string `yaml:"name"`
		Filters []struct {
			Name  string
			Value string
		} `yaml:"filters"`
		Columns []string `yaml:"columns"`
		Ignore  []string `yaml:"ignore"`
	} `yaml:"tables"`
}

func DigestSchema(fpath string) (Schema, error) {
	schema := Schema{}
	yml, err := os.ReadFile(fpath)

	if err != nil {
		return schema, err
	}

	if err := yaml.Unmarshal([]byte(yml), &schema); err != nil {
		return schema, err
	}

	return schema, nil
}
