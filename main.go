package main

import (
	"final-project-pbi-btpns/database"
	"final-project-pbi-btpns/router"
)

func main() {
	// koneksi ke db
	database.ConnectToDb()

	//migrasi database
	database.SyncDatabase()

	// setup router
	r := router.SetupRouter()

	// run server
	r.Run(":8080")

}
