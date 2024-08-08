package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/backend/db/sqlite"
	"social-network/backend/routes"
)

func main() {
	dbPath := "./db/database.db"
	migrationsPath := "./db/migrations/sqlite" 

	db, err := sqlite.ConnectAndMigrateDb(dbPath, migrationsPath)
	if err != nil {
		log.Fatalf("Failed to connect or migrate the database: %v", err)
	}
	defer db.Close()

	mux := routes.SetupRoutes()
	fmt.Println("Starting server on: http://localhost:8080\nCtrl+c for exit")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
