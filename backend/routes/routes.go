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
	// mux.Handle("/testLoggedIn", middleware.RequireLogin(http.HandlerFunc(middleware.DummyCheck)))
	mux.HandleFunc("/create-posts", handlers.CreatePost)
	mux.HandleFunc("/posts", handlers.GetPosts)
	mux.HandleFunc("/create-comment", handlers.CreateComment)
	mux.HandleFunc("/comments", handlers.GetComments)
	mux.HandleFunc("/search", handlers.SearchUser)
	mux.HandleFunc("/followers", HandleFollowers)
	mux.HandleFunc("/followers/status", HandleFollowers)
	mux.HandleFunc("/followers/pending", HandleFollowers)
	mux.HandleFunc("/followers/accept", HandleFollowers)
	mux.HandleFunc("/followers/reject", HandleFollowers)
	mux.HandleFunc("/followers/followersList", HandleFollowers)
	mux.HandleFunc("/groups", HandleGroups)
	mux.HandleFunc("/groups/join-request", HandleGroups)
	mux.HandleFunc("/groups/invite", HandleGroups)
	mux.HandleFunc("/groups/members", HandleGroups)
	mux.HandleFunc("/groups/handle-request", HandleGroups)
	mux.HandleFunc("/groups/handle-invites", HandleGroups)
	mux.HandleFunc("/groups/pending-requests", HandleGroups)
	mux.HandleFunc("/groups/pending-invites", HandleGroups)
	mux.HandleFunc("/groups/get-users", HandleGroups)
	mux.HandleFunc("/groups/events", HandleGroupEvents)
	mux.HandleFunc("/groups/events-react", HandleGroupEvents)
	mux.HandleFunc("/groups/members-with-reactions", HandleGroupEvents)

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
	if r.Method == http.MethodGet && r.URL.Path == "/followers/followersList" {
		handlers.GetFollowLists(w, r)
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
	if r.Method == http.MethodGet && r.URL.Path == "/groups/get-users" {
		handlers.GetUserGroups(w, r)
		return
	}

	if r.Method == http.MethodGet && r.URL.Path == "/groups/pending-invites" {
		handlers.GetPendingGroupInvites(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/groups/members" {
		handlers.GetGroupMembers(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/groups/pending-requests" {
		handlers.GetPendingGroupJoin(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/groups/join-request" {
		handlers.JoinGroupRequest(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/groups/invite" {
		handlers.InviteToGroup(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/groups/handle-invites" {
		handlers.InviteRequestHandler(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/groups/handle-request" {
		handlers.JoinRequestHandler(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		handlers.CreateGroup(w, r)
	case http.MethodGet:
		handlers.GetGroupsWithMembers(w, r)
	default:
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func HandleGroupEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/groups/members-with-reactions" {
		handlers.GetGroupMembersWithReactions(w, r)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/groups/events-react" {
		handlers.EventReaction(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		handlers.CreateEvent(w, r)
	case http.MethodGet:
		handlers.GetEvents(w, r)
	default:
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
