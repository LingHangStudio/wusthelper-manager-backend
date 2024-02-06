package database

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	PrintSql bool
	Pool     struct {
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime int
	}
	ExtraParams map[string]string
}
