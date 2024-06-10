package config

import (
	"challenge-client-server/Entities"
	"challenge-client-server/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	var err error
	Db, err = gorm.Open(sqlite.Open("client-server.db"), &gorm.Config{})
	utils.HandlerError(err)

	err = Db.AutoMigrate(&Entities.Quotation{})
	utils.HandlerError(err)
}
