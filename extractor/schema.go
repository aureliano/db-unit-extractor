package extractor

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"gopkg.in/yaml.v2"
)

const NameMaxLength = 80

var ErrSchemaValidation = errors.New("validation")
var nameRegExp = regexp.MustCompile(`^[a-zA-Z_]\w+$`)

type ConverterSchema string
type ColumnSchema string
type IgnoreSchema string

type FilterSchema struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type TableSchema struct {
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

func DigestSchema(fpath string) (Schema, error) {
	schema := Schema{}
	yml, err := os.ReadFile(fpath)

	if err != nil {
		return schema, err
	}

	if err = yaml.Unmarshal(yml, &schema); err != nil {
		return schema, err
	}

	return schema, nil
}

func (s Schema) Validate() error {
	if len(s.Converters) > 0 {
		if err := validateConverters(s.Converters); err != nil {
			return err
		}
	}

	return validateTables(s.Tables)
}

func (c ConverterSchema) Validate() error {
	if !dataconv.ConverterExists(string(c)) {
		return fmt.Errorf("%w: converter '%s' not found", ErrSchemaValidation, c)
	}

	return nil
}

func (t TableSchema) Validate() error {
	if err := validateName(t.Name); err != nil {
		return err
	}

	for _, filter := range t.Filters {
		if err := filter.Validate(); err != nil {
			return fmt.Errorf("table '%s' %w", t.Name, err)
		}
	}

	for _, column := range t.Columns {
		if err := column.Validate(); err != nil {
			return fmt.Errorf("table '%s' %w", t.Name, err)
		}
	}

	for _, column := range t.Ignore {
		if err := column.Validate(); err != nil {
			return fmt.Errorf("table '%s' %w", t.Name, err)
		}
	}

	if len(t.Columns) > 0 && len(t.Ignore) > 0 {
		return fmt.Errorf("%w: table '%s' with columns and ignore set (excludents)", ErrSchemaValidation, t.Name)
	}

	return nil
}

func (f FilterSchema) Validate() error {
	if len(f.Value) == 0 {
		return fmt.Errorf("%w: empty filter value '%s'", ErrSchemaValidation, f.Name)
	}

	return validateName(f.Name)
}

func (c ColumnSchema) Validate() error {
	return validateName(string(c))
}

func (c IgnoreSchema) Validate() error {
	return validateName(string(c))
}

func validateConverters(converters []ConverterSchema) error {
	for _, converter := range converters {
		if err := converter.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func validateTables(tables []TableSchema) error {
	if len(tables) == 0 {
		return fmt.Errorf("%w: no table provided", ErrSchemaValidation)
	}

	for _, table := range tables {
		if err := table.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func validateName(name string) error {
	if len(name) > NameMaxLength {
		return fmt.Errorf("%w: '%s' invalid name size", ErrSchemaValidation, name)
	}

	if !nameRegExp.MatchString(name) {
		return fmt.Errorf("%w: '%s' invalid name", ErrSchemaValidation, name)
	}

	return nil
}
