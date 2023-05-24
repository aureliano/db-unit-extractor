package extractor

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
)

type Conf struct {
	SchemaPath string
	reader.DataSource
	References map[string]interface{}
}

type dbResponse struct {
	table string
	data  []map[string]interface{}
	err   error
}

var (
	ErrExtractor      = errors.New("extractor")
	filterValueRegExp = regexp.MustCompile(`^\$\{([^\}]+)\}$`)
)

func Extract(conf Conf, db reader.DBReader) error {
	schema, err := schema.DigestSchema(conf.SchemaPath)
	if err != nil {
		return err
	}

	if db == (reader.DBReader)(nil) {
		db, err = reader.NewReader(conf.DataSource, nil)
		if err != nil {
			return err
		}
	}

	for k, v := range conf.References {
		schema.Refs[k] = v
	}

	return extract(schema, db)
}

func extract(model schema.Model, db reader.DBReader) error {
	for _, tables := range model.GroupedTables() {
		respChan := make(chan dbResponse)
		responses := make([]dbResponse, 0, len(tables))
		tbSize := len(tables)

		for _, table := range tables {
			filters, err := resolveTableFilters(table, model.Refs)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrExtractor, err)
			}

			go fetchData(respChan, table, db, convToStr(model.Converters), filters)
		}

		counter := 0
		for res := range respChan {
			if res.err != nil {
				return fmt.Errorf("%w: %w", ErrExtractor, res.err)
			}

			responses = append(responses, res)
			counter++

			if counter >= tbSize {
				close(respChan)
			}
		}

		updateReferences(model, responses)
	}

	return nil
}

func fetchData(c chan dbResponse, table schema.Table,
	db reader.DBReader, converters []string, filters [][]interface{}) {
	columns, err := db.FetchColumnsMetadata(table)
	if err != nil {
		c <- dbResponse{err: err}
		return
	}

	data, err := db.FetchData(table.Name, columns, converters, filters)
	if err != nil {
		c <- dbResponse{err: err}
		return
	}

	c <- dbResponse{
		table: table.Name,
		data:  data,
		err:   err,
	}
}

func updateReferences(model schema.Model, responses []dbResponse) {
	for _, response := range responses {
		for _, record := range response.data {
			for k, v := range record {
				key := fmt.Sprintf("%s.%s", response.table, k)
				if _, exist := model.Refs[key]; exist {
					model.Refs[key] = v
				}
			}
		}
	}
}

func resolveTableFilters(table schema.Table, references map[string]interface{}) ([][]interface{}, error) {
	size := len(table.Filters)
	filters := make([][]interface{}, size)
	const pair = 2

	for i := range filters {
		filters[i] = make([]interface{}, pair)
	}

	for i := 0; i < size; i++ {
		filter := table.Filters[i]
		var value interface{}

		matches := filterValueRegExp.FindAllStringSubmatch(filter.Value, -1)
		if matches != nil {
			key := matches[0][1]

			if v, exists := references[key]; exists {
				value = v
			} else {
				return nil, fmt.Errorf("filter %s.%s not found '%s'", table.Name, filter.Name, filter.Value)
			}
		} else {
			value = filter.Value
		}

		filters[i][0] = filter.Name
		filters[i][1] = value
	}

	return filters, nil
}

func convToStr(conv []schema.Converter) []string {
	res := make([]string, len(conv))
	for i, c := range conv {
		res[i] = string(c)
	}

	return res
}
