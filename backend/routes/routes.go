package routes

import (
	"net/http"
	"social-network/handlers"
	"social-network/middleware"
	"social-network/wsocket"
)

func SetupRoutes() *http.ServeMux {
	manager := wsocket.NewManager()
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/ws", manager.Serve_WS)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/getAvatarPath", handlers.GetAvatarPath)
	mux.Handle("/getUserData", middleware.RequireLogin(http.HandlerFunc(handlers.GetUserData)))
	mux.Handle("/testLoggedIn", middleware.RequireLogin(http.HandlerFunc(middleware.DummyCheck)))
	return mux
}
