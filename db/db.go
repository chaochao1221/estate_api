package db

import (
	"estate/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var Db *xorm.Engine

func Init() {
	dbconn, err := utils.LoadDbConfig("estate")
	if err == nil {
		Db, err = xorm.NewEngine("mysql", dbconn)
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		Db.SetMaxOpenConns(1000)
		Db.SetMaxIdleConns(300)
	}
}

func Close() {
	Db.Close()
}
