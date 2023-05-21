package extractor

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

const NameMaxLength = 80

var (
	ErrSchemaFile          = errors.New("schema-file")
	ErrSchemaValidation    = errors.New("validation")
	ErrTableClassification = errors.New("classification")
	nameRegExp             = regexp.MustCompile(`^[a-zA-Z_]\w+$`)
	filterReferenceRegExp  = regexp.MustCompile(`^\$\{(\w+)\.(\w+)\}$`)
)

type ConverterSchema string
type ColumnSchema string
type IgnoreSchema string

type FilterSchema struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type TableSchema struct {
	GroupID int
	Name    string         `yaml:"name"`
	Filters []FilterSchema `yaml:"filters"`
	Columns []ColumnSchema `yaml:"columns"`
	Ignore  []IgnoreSchema `yaml:"ignore"`
}

type Schema struct {
	Refs       map[string]interface{}
	Converters []ConverterSchema `yaml:"converters"`
	Tables     []TableSchema     `yaml:"tables"`
}

type SchemaValidator interface {
	Validate() error
}

type SchemaClassifier interface {
	Classify() error
	GroupedTables() [][]TableSchema
}

func DigestSchema(fpath string) (Schema, error) {
	schema := Schema{}
	yml, err := os.ReadFile(fpath)

	if err != nil {
		return schema, fmt.Errorf("%w: %w", ErrSchemaFile, err)
	}

	if err = yaml.UnmarshalStrict(yml, &schema); err != nil {
		return schema, fmt.Errorf("%w: %w", ErrSchemaFile, err)
	}

	if err = schema.Validate(); err != nil {
		return schema, err
	}

	schema.Refs = fetchReferences(schema)

	return schema, schema.Classify()
}

func fetchReferences(s Schema) map[string]interface{} {
	refs := make(map[string]interface{})

	for _, table := range s.Tables {
		for _, filter := range table.Filters {
			matches := filterReferenceRegExp.FindAllStringSubmatch(filter.Value, -1)
			if matches != nil {
				key := fmt.Sprintf("%s.%s", matches[0][1], matches[0][2])
				refs[key] = nil
			}
		}
	}

	return refs
}
