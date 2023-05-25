package reader

import (
	"fmt"
	"strings"
)

type DataSource struct {
	DBMSName    string
	Username    string
	Password    string
	Database    string
	Hostname    string
	Port        int
	MaxOpenConn int
	MaxIdleConn int
}

type DSN interface {
	DSName() string
}

const dsnTemplate = "%s://%s:%s@%s:%d/%s"

func NewDataSource() DataSource {
	return DataSource{
		MaxOpenConn: 1,
		MaxIdleConn: 1,
	}
}

func (ds DataSource) DSName() string {
	if ds.DBMSName == "sqlite3" {
		return "file:test.db?cache=shared&mode=memory"
	}

	return fmt.Sprintf(
		dsnTemplate,
		strings.ToLower(ds.DBMSName),
		ds.Username,
		ds.Password,
		ds.Hostname,
		ds.Port,
		ds.Database,
	)
}
