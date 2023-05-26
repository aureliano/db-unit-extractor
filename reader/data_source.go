package reader

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
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
	DB          *sql.DB
}

type DSN interface {
	DSName() string
}

type DBConnector interface {
	Connect() error
	IsConnected() bool
}

const (
	dsnTemplate  = "%s://%s:%s@%s:%d/%s"
	MaxDBTimeout = time.Second * 30
)

func NewDataSource() *DataSource {
	return &DataSource{
		MaxOpenConn: 1,
		MaxIdleConn: 1,
	}
}

func (ds *DataSource) DSName() string {
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

func (ds *DataSource) Connect(timeout time.Duration) error {
	if ds.DB != nil {
		return nil
	}

	db, err := sql.Open(ds.DBMSName, ds.DSName())
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(ds.MaxOpenConn)
	db.SetMaxIdleConns(ds.MaxIdleConn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	ds.DB = db

	return nil
}

func (ds *DataSource) IsConnected() bool {
	return ds.DB != nil
}
