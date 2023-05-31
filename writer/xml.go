package writer

type XMLWriter struct{}

func (XMLWriter) WriteHeader() {
	_ = 0
}

func (XMLWriter) WriteFooter() {
	_ = 0
}

func (XMLWriter) Write(table string, rows []map[string]interface{}) error {
	return nil
}
