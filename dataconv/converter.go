package dataconv

var converterIds = make(map[string]Converter)

type Converter interface {
	Convert(interface{}) (interface{}, error)
}

func ConverterExists(id string) bool {
	return converterIds[id] != nil
}

func RegisterConverter(id string, converter Converter) {
	converterIds[id] = converter
}

func RegisterConverters() {
	RegisterConverter("date-time-iso8601", DateTimeISO8601Converter{})
	RegisterConverter("blob", BlobConverter{})
}
