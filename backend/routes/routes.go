package routes

import (
	"net/http"
	"social-network/handlers"
	"social-network/middleware"
	"social-network/sockets"
)

func SetupRoutes() *http.ServeMux {

	manager := sockets.NewManager()
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
	mux.HandleFunc("/followers", HandleFollowers)
	mux.HandleFunc("/followers/status", HandleFollowers)
	mux.HandleFunc("/followers/pending", HandleFollowers)
	mux.HandleFunc("/followers/accept", HandleFollowers)
	mux.HandleFunc("/followers/reject", HandleFollowers)
	mux.HandleFunc("/groups", HandleGroups)
	mux.HandleFunc("/groups/join", HandleGroups)

	return mux
}
func HandleFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/followers/accept" {
		handlers.AcceptFollowRequest(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/followers/reject" {
		handlers.RejectFollowRequest(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/followers/status" {
		handlers.GetFollowStatus(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/followers/pending" {
		handlers.GetPendingFollowRequests(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		handlers.FollowUser(w, r)
	case http.MethodDelete:
		handlers.UnfollowUser(w, r)
	default:
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func HandleGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/groups/join" {
		handlers.JoinGroupHandler(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		handlers.CreateGroup(w, r)
	case http.MethodGet:
		handlers.GetAllGroups(w, r)
	default:
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
