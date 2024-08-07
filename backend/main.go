package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/backend/routes"
)

func main() {
	mux := routes.SetupRoutes()
	fmt.Println("Starting server on: http://localhost:8080\nCtrl+c for exit")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
