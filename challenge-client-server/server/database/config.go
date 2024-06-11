package database

import (
	"challenge-server/entities"
	"challenge-server/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	var err error
	Db, err = gorm.Open(sqlite.Open("client-server.db"), &gorm.Config{})
	utils.HandlerError(err)

	err = Db.AutoMigrate(&entities.Quotation{})
	utils.HandlerError(err)
}
