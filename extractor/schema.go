package extractor

import (
	"errors"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

const NameMaxLength = 80

var (
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
	Converters []ConverterSchema `yaml:"converters"`
	Tables     []TableSchema     `yaml:"tables"`
}

type SchemaValidator interface {
	Validate() error
}

type SchemaClassifier interface {
	Classify() error
}

func DigestSchema(fpath string) (Schema, error) {
	schema := Schema{}
	yml, err := os.ReadFile(fpath)

	if err != nil {
		return schema, err
	}

	if err = yaml.UnmarshalStrict(yml, &schema); err != nil {
		return schema, err
	}

	if err = schema.Validate(); err != nil {
		return schema, err
	}

	return schema, schema.Classify()
}
