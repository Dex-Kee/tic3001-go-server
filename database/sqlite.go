package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initSqlite() *gorm.DB {
	sqliteConn, err := gorm.Open(sqlite.Open("server.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	s, err := sqliteConn.DB()
	if err != nil {
		log.Panic(err)
	}

	if err = s.Ping(); err != nil {
		log.Panic(err)
	}

	return sqliteConn
}
