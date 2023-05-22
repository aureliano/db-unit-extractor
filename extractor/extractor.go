package extractor

import (
	"fmt"
	"sync"

	"github.com/aureliano/db-unit-extractor/reader"
	"github.com/aureliano/db-unit-extractor/schema"
)

type Conf struct {
	SchemaPath string
	reader.DataSource
}

func Extract(conf Conf) error {
	schema, err := schema.DigestSchema(conf.SchemaPath)
	if err != nil {
		return err
	}

	db, err := reader.NewReader(conf.DataSource)
	if err != nil {
		return err
	}

	return extract(conf.DataSource, schema, db)
}

func extract(ds reader.DataSource, schema schema.Model, db reader.DBReader) error {
	groupedTables := schema.GroupedTables()
	var wg sync.WaitGroup

	for _, tables := range groupedTables {
		for _, table := range tables {
			fmt.Println(table)
		}

		/*wg.Add(1)
		go func(table TableSchema) {
			defer wg.Done()

			db.FetchColumnsMetadata(table.Name, table.Columns, table.Ignore)
		}(table)*/
	}

	wg.Wait()

	return nil
}
