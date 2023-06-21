package dataconv

type Converter interface {
	Convert(string, interface{}) (interface{}, error)
	Handle(interface{}) bool
}

const (
	DateTimeISO8601ID = "date-time-iso8601"
	BlobConverterID   = "blob"
)

var converterIds = make(map[string]Converter)

func ConverterExists(id string) bool {
	return converterIds[id] != nil
}

func RegisterConverter(id string, converter Converter) {
	converterIds[id] = converter
}

func RegisterConverters() {
	RegisterConverter(DateTimeISO8601ID, DateTimeISO8601Converter{})
	RegisterConverter(BlobConverterID, BlobConverter{})
}

func GetConverter(id string) Converter {
	return converterIds[id]
}
