package database

import (
	"final-project-pbi-btpns/models"
	"log"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Photo{})
	log.Println("database migration completed")
}
