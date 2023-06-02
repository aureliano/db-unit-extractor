package schema

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

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

type Converter string
type Column string
type Ignore string

type Filter struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Table struct {
	GroupID int
	Name    string   `yaml:"name"`
	Filters []Filter `yaml:"filters"`
	Columns []Column `yaml:"columns"`
	Ignore  []Ignore `yaml:"ignore"`
}

type Model struct {
	Refs       map[string]interface{}
	Converters []Converter `yaml:"converters"`
	Tables     []Table     `yaml:"tables"`
}

type Validator interface {
	Validate() error
}

type Classifier interface {
	Classify() error
	GroupedTables() [][]Table
}

type DataTable interface {
	SelectColumns() []string
	FormattedSelectColumns() string
}

func DigestSchema(fpath string) (Model, error) {
	schema := Model{}
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

func (t Table) SelectColumns() []string {
	var columns []string
	const wildCardFrom = 1

	switch {
	case len(t.Columns) > 0:
		columns = make([]string, len(t.Columns))
		for i, c := range t.Columns {
			columns[i] = string(c)
		}
	case len(t.Ignore) > 0:
		columns = make([]string, len(t.Ignore))
		for i, c := range t.Ignore {
			columns[i] = string(c)
		}
	default:
		columns = make([]string, wildCardFrom)
		columns[0] = "*"
	}

	return columns
}

func (t Table) FormattedSelectColumns() string {
	return fmt.Sprintf("'%s'", strings.Join(t.SelectColumns(), "', '"))
}

func fetchReferences(s Model) map[string]interface{} {
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
