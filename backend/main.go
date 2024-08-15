package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/backend/db/sqlite"
	"social-network/backend/middleware"
	"social-network/backend/routes"
)

func main() {

	migrationsPath := "./db/migrations/sqlite"

	_, err := sqlite.ConnectAndMigrateDb(migrationsPath)
	if err != nil {
		log.Fatalf("Failed to connect or migrate the database: %v", err)
	}
	defer sqlite.Db.Close()

	mux := routes.SetupRoutes()
	fmt.Println("Starting server on: http://localhost:8080\nCtrl+c for exit")
	log.Fatal(http.ListenAndServe(":8080", middleware.CorsMiddleWare(mux)))
}
