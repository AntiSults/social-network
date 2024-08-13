package routes

import (
	"net/http"
	"social-network/backend/handlers"
	"social-network/backend/middleware"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/getAvatarPath", handlers.GetAvatarPath)
	mux.Handle("/testLoggedIn", middleware.RequireLogin(http.HandlerFunc(middleware.DummyCheck)))
	return mux
}
