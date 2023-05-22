package extractor

import (
	"errors"
	"fmt"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
)

type Conf struct {
	SchemaPath string
	reader.DataSource
}

type dbResponse struct {
	table   string
	filters [][]interface{}
	data    []map[string]interface{}
	err     error
}

var ErrExtractor = errors.New("extractor")

func Extract(conf Conf, db reader.DBReader) error {
	schema, err := schema.DigestSchema(conf.SchemaPath)
	if err != nil {
		return err
	}

	if db == (reader.DBReader)(nil) {
		db, err = reader.NewReader(conf.DataSource)
		if err != nil {
			return err
		}
	}

	return extract(schema, db)
}

func extract(model schema.Model, db reader.DBReader) error {
	references := prepareReferences(model)
	respChan := make(chan dbResponse)

	defer close(respChan)

	for _, tables := range model.GroupedTables() {
		responses := make([]dbResponse, 0, len(tables))

		for _, table := range tables {
			filters := resolveTableFilters(table, references)

			go fetchData(respChan, table, db, convToStr(model.Converters), filters)
		}

		for res := range respChan {
			if res.err != nil {
				return fmt.Errorf("%w: %w", ErrExtractor, res.err)
			}

			responses = append(responses, res)
		}
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
		table:   table.Name,
		filters: updateReferences(data),
		data:    data,
		err:     err,
	}
}

func updateReferences(_ []map[string]interface{}) [][]interface{} {
	panic("unimplemented")
}

func resolveTableFilters(_ schema.Table, _ map[string]interface{}) [][]interface{} {
	panic("unimplemented")
}

func prepareReferences(_ schema.Model) map[string]interface{} {
	panic("unimplemented")
}

func convToStr(conv []schema.Converter) []string {
	res := make([]string, len(conv))
	for i, c := range conv {
		res[i] = string(c)
	}

	return res
}
