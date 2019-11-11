package model

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

var SqliteConn *gorm.DB

func init() {
	// init sqlite database
	db, err := gorm.Open("sqlite3", "data/factory.db")
	if err != nil {
		log.Println(errors.WithMessage(err, "can not open factory database(with sqlite3)"))
	}
	db.AutoMigrate(&Factory{})
	db.AutoMigrate(&Instance{})
	db.AutoMigrate(&ServiceRegistry{})
	SqliteConn = db
}
