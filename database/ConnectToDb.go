package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbError error

func ConnectToDb() {
	connectionString := "root:@tcp(localhost:3306)/pbi_api?parseTime=true"

	database, dbError := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if dbError != nil {
		log.Fatal(dbError)
		panic("failed to connect database")
	}

	log.Println("success to connection to database")

	DB = database
}
