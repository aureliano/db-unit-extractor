package schema

import (
	"fmt"

	"github.com/aureliano/db-unit-extractor/dataconv"
)

func (s Model) Validate() error {
	if len(s.Converters) > 0 {
		if err := validateConverters(s.Converters); err != nil {
			return err
		}
	}

	return validateTables(s.Tables)
}

func (c Converter) Validate() error {
	if !dataconv.ConverterExists(string(c)) {
		return fmt.Errorf("%w: converter '%s' not found", ErrSchemaValidation, c)
	}

	return nil
}

func (t Table) Validate() error {
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

	cols := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		cols[i] = string(c)
	}

	c := repeatedValue(cols)
	if c != "" {
		return fmt.Errorf("%w: repeated column '%s' in table '%s", ErrSchemaValidation, c, t.Name)
	}

	return nil
}

func (f Filter) Validate() error {
	if len(f.Value) == 0 {
		return fmt.Errorf("%w: empty filter value '%s'", ErrSchemaValidation, f.Name)
	}

	return validateName(f.Name)
}

func (c Column) Validate() error {
	return validateName(string(c))
}

func (c Ignore) Validate() error {
	return validateName(string(c))
}

func validateTables(tables []Table) error {
	if len(tables) == 0 {
		return fmt.Errorf("%w: no table provided", ErrSchemaValidation)
	}

	for _, table := range tables {
		if err := table.Validate(); err != nil {
			return err
		}
	}

	tbls := make([]string, len(tables))
	for i, t := range tables {
		tbls[i] = t.Name
	}

	tb := repeatedValue(tbls)
	if tb != "" {
		return fmt.Errorf("%w: repeated table '%s'", ErrSchemaValidation, tb)
	}

	return nil
}

func validateConverters(converters []Converter) error {
	for _, converter := range converters {
		if err := converter.Validate(); err != nil {
			return err
		}
	}

	convs := make([]string, len(converters))
	for i, c := range converters {
		convs[i] = string(c)
	}

	conv := repeatedValue(convs)
	if conv != "" {
		return fmt.Errorf("%w: repeated converter '%s'", ErrSchemaValidation, conv)
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

func repeatedValue(values []string) string {
	for i, value := range values {
		for j, str := range values {
			if i != j && value == str {
				return value
			}
		}
	}

	return ""
}
