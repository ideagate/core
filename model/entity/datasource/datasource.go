package datasource

import (
	"github.com/bayu-aditya/ideagate/backend/core/model/constant"
	"gorm.io/gorm"
)

// DataSource entity for table `datasource`
type DataSource struct {
	Id     string
	Type   constant.DataSourceType // Ex: mysql, postgresql, redis, rest, etc
	Config Config

	// connection
	MysqlConn *gorm.DB `json:"-"`
}

// Config entity for json struct DataSource.Config
type Config struct {
	Host     string `json:",omitempty"`
	DB       string `json:",omitempty"`
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
}
