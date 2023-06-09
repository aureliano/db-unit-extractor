package dataconv

var converterIds = make(map[string]Converter)

type Converter interface {
	Convert(source interface{}) (target interface{})
}

func ConverterExists(id string) bool {
	return converterIds[id] != nil
}

func RegisterConverter(id string, converter Converter) {
	converterIds[id] = converter
}
