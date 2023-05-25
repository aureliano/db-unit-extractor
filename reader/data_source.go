package reader

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

const dsnTemplate = "%s://%s:%s@%s:%d/%s"

func NewDataSource() DataSource {
	return DataSource{
		MaxOpenConn: 1,
		MaxIdleConn: 1,
	}
}
