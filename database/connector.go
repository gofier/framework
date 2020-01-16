package database

import (
	"github.com/gofier/framework/config"
	"github.com/gofier/framework/database/driver"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var dber IDatabase

func Init() {
	dber, db = setConnection("default")
}

func setConnection(conn string) (dber IDatabase, _db *gorm.DB) {
	_conn := conn

	if conn == "default" {
		_conn = config.GetString("database." + conn)
		if _conn == "" {
			panic("database connection parse error")
		}
	}

	switch _conn {
	case "mysql":
		dber = driver.NewMysql(_conn)
		break
	default:
		panic("incorrect database connection provided")
	}

	_db, err := gorm.Open(_conn, dber.ConnectionArgs())
	if err != nil {
		panic("failed to connect database")
	}

	err = _db.DB().Ping()
	if err != nil {
		panic("failed to connect database by ping")
	}

	if config.GetBoolean("app.debug") {
		_db = _db.Debug().LogMode(true)
	}

	_db.DB().SetMaxIdleConns(config.GetInt("database.max_idle_connections"))
	_db.DB().SetMaxOpenConns(config.GetInt("database.max_open_connections"))

	return dber, _db
}

func Connection(conn string) (_db *gorm.DB) {
	_, _db = setConnection(conn)
	return _db
}

func DB() *gorm.DB {
	return db
}

func Prefix() string {
	return dber.Prefix()
}
