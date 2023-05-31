package reader

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type DataSource struct {
	DSN         string
	MaxOpenConn int
	MaxIdleConn int
	DB          *sql.DB
}

type DBConnector interface {
	Connect() error
	IsConnected() bool
	DriverName() string
	ConnectionURL() string
}

const MaxDBTimeout = time.Second * 30

func NewDataSource() *DataSource {
	return &DataSource{
		MaxOpenConn: 1,
		MaxIdleConn: 1,
	}
}

func (ds *DataSource) Connect(timeout time.Duration) error {
	if ds.DB != nil {
		return nil
	}

	db, err := sql.Open(ds.DriverName(), ds.ConnectionURL())
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

func (ds *DataSource) DriverName() string {
	if ds.DSN == "" {
		return ""
	}

	index := strings.Index(ds.DSN, "://")
	return ds.DSN[:index]
}

func (ds *DataSource) ConnectionURL() string {
	if ds.DSN == "" {
		return ""
	}

	const sqlite3Prefix = "sqlite3://"
	if strings.HasPrefix(ds.DSN, sqlite3Prefix) {
		return ds.DSN[len(sqlite3Prefix):]
	}

	return ds.DSN
}
