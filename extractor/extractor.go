package extractor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aureliano/db-unit-extractor/dataconv"
	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
	"github.com/aureliano/db-unit-extractor/writer"
)

type Conf struct {
	SchemaPath      string
	DSN             string
	MaxOpenConn     int
	MaxIdleConn     int
	OutputTypes     []string
	FormattedOutput bool
	OutputDir       string
	References      map[string]interface{}
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

func Extract(conf Conf, db reader.DBReader, writers []writer.FileWriter) error {
	dataconv.RegisterConverters()
	schema, err := schema.DigestSchema(conf.SchemaPath)
	if err != nil {
		return err
	}

	if db == (reader.DBReader)(nil) {
		ds := reader.NewDataSource()
		ds.DSN = conf.DSN
		ds.MaxOpenConn = conf.MaxOpenConn
		ds.MaxIdleConn = conf.MaxIdleConn

		db, err = reader.NewReader(ds)
		if err != nil {
			return err
		}
	}

	if len(writers) == 0 {
		writers = make([]writer.FileWriter, len(conf.OutputTypes))

		for i, outputTp := range conf.OutputTypes {
			fname := filepath.Base(conf.SchemaPath)
			fname = fname[:strings.LastIndex(fname, ".")]

			fc := writer.FileConf{
				Type: outputTp, Formatted: conf.FormattedOutput, Directory: conf.OutputDir, Name: fname,
			}
			fw, e := writer.NewWriter(fc)
			if e != nil {
				return e
			}
			writers[i] = fw
		}
	}

	for k, v := range conf.References {
		schema.Refs[strings.ToLower(k)] = v
	}

	return extract(schema, db, writers)
}

func extract(model schema.Model, db reader.DBReader, writers []writer.FileWriter) error {
	cw := launchWriters(writers)

	if err := launchReaders(model, db, cw); err != nil {
		return err
	}

	for i := 0; i < len(cw); i++ {
		cw[i] <- dbResponse{}
		_ = writers[i].WriteFooter()
	}

	return nil
}

func launchWriters(writers []writer.FileWriter) []chan dbResponse {
	chanWriters := make([]chan dbResponse, len(writers))

	for i, w := range writers {
		chanWriters[i] = make(chan dbResponse)
		go writeData(chanWriters[i], w)
	}

	return chanWriters
}

func launchReaders(model schema.Model, db reader.DBReader, writers []chan dbResponse) error {
	converters := make([]dataconv.Converter, 0)
	for _, id := range model.Converters {
		converters = append(converters, dataconv.GetConverter(string(id)))
	}

	for _, tables := range model.GroupedTables() {
		respChan := make(chan dbResponse)
		tbSize := len(tables)

		for _, table := range tables {
			filters, err := resolveTableFilters(table, model.Refs)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrExtractor, err)
			}

			go fetchData(respChan, table, db, converters, filters)
		}

		counter := 0
		for res := range respChan {
			if res.err != nil {
				return fmt.Errorf("%w: %w", ErrExtractor, res.err)
			}

			updateReferences(model, res)
			counter++

			if counter >= tbSize {
				close(respChan)
			}

			for _, w := range writers {
				w <- res
			}
		}
	}

	return nil
}

func fetchData(c chan dbResponse, table schema.Table,
	db reader.DBReader, converters []dataconv.Converter, filters [][]interface{}) {
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

func writeData(c chan dbResponse, w writer.FileWriter) {
	if err := w.WriteHeader(); err != nil {
		shutdown(err)
	}

	for res := range c {
		if res.data != nil {
			if err := w.Write(res.table, res.data); err != nil {
				shutdown(err)
			}
		}
	}
}

func updateReferences(model schema.Model, response dbResponse) {
	for _, record := range response.data {
		for k, v := range record {
			key := strings.ToLower(fmt.Sprintf("%s.%s", response.table, k))
			if _, exist := model.Refs[key]; exist {
				model.Refs[key] = v
			} else {
				key = strings.ToLower(fmt.Sprintf("%s.%s[@]", response.table, k))
				if _, exist = model.Refs[key]; exist {
					if model.Refs[key] == nil {
						model.Refs[key] = make([]interface{}, 0)
					}

					model.Refs[key] = append(model.Refs[key].([]interface{}), v)
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
			key := strings.ToLower(matches[0][1])

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

func shutdown(err error) {
	fmt.Fprintf(os.Stdout, "%s: %s\n", ErrExtractor.Error(), err.Error())
	os.Exit(1)
}
