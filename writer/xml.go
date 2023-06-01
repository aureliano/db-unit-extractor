package writer

type XMLWriter struct{}

func (XMLWriter) WriteHeader() error {
	return nil
}

func (XMLWriter) WriteFooter() error {
	return nil
}

func (XMLWriter) Write(table string, rows []map[string]interface{}) error {
	return nil
}
