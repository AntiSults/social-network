package routes

import (
	"social-network/backend/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
	return router
}