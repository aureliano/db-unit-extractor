package extractor

import "github.com/aureliano/db-unit-extractor/reader"

type Conf struct {
	SchemaPath string
	reader.DataSource
}

func Extract(conf Conf) error {
	schema, err := DigestSchema(conf.SchemaPath)
	if err != nil {
		return err
	}

	db, err := reader.NewReader(conf.DataSource)
	if err != nil {
		return err
	}

	return extract(schema, db)
}

func extract(schema Schema, db reader.DBReader) error {
	return nil
}
