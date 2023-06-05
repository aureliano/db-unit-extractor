package main

import (
	"encoding/xml"
	"os"
)

type Tag struct {
	Table  string
	Fields []Field
}

type XMLDataSet struct {
	XMLName xml.Name `xml:"dataset"`
	Records []Tag    `xml:",any"`
}

func (r *Tag) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	r.Table = start.Name.Local
	fields := make([]Field, len(start.Attr))
	r.Fields = fields

	for i, attr := range start.Attr {
		fields[i] = Field{Name: attr.Name.Local, Value: attr.Value}
	}

	return d.DecodeElement(&r.Fields, &start)
}

func parseXML(path string) (*DataSet, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ds XMLDataSet
	if err = xml.Unmarshal(byteValue, &ds); err != nil {
		return nil, err
	}

	dataset := DataSet{
		Tables: make([]Table, len(ds.Records)),
	}

	for i, record := range ds.Records {
		fields := make([]Field, len(record.Fields)-1)

		for j, field := range record.Fields {
			if field.Name != "" {
				fields[j] = field
			}
		}

		dataset.Tables[i] = Table{
			Name:   record.Table,
			Fields: fields,
		}
	}

	return &dataset, nil
}
