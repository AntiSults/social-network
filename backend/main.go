package main

import (
	"log"
	"net/http"
	"social-network/backend/routes"
)

func main(){
	router := routes.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}