package routes

import (
	"net/http"
	"social-network/handlers"
	"social-network/middleware"
	"social-network/sockets"
)

func SetupRoutes() *http.ServeMux {

	manager := sockets.NewManager() //need it here as instance of Manager struct (Serve_WS is method)
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/ws", manager.Serve_WS)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/getAvatarPath", handlers.GetAvatarPath)
	mux.Handle("/getUserData", middleware.RequireLogin(http.HandlerFunc(handlers.GetUserData)))
	mux.Handle("/testLoggedIn", middleware.RequireLogin(http.HandlerFunc(middleware.DummyCheck)))
	mux.HandleFunc("/create-posts", handlers.CreatePost)
	mux.HandleFunc("/posts", handlers.GetPosts)
	mux.HandleFunc("/search", handlers.SearchUser)
	return mux
}
