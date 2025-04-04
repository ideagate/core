package constant

type DataSourceType string

var (
	DataSourceTypeMysql      DataSourceType = "mysql"
	DataSourceTypePostgresql DataSourceType = "postgresql"
	DataSourceTypeRedis      DataSourceType = "redis"
	DataSourceTypeRest       DataSourceType = "rest"
)
