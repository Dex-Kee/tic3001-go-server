package database

import "gorm.io/gorm"

var Conn = new(config)

type config struct {
	connectionIns *gorm.DB
}

func InitDatabase() {
	sqlite := initSqlite()
	Conn.connectionIns = sqlite
}

func (c config) GetConnection() *gorm.DB {
	return c.connectionIns
}
