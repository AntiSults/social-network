package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/structs"
	"strconv"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID     int `json:"userId"`
		FollowerID int `json:"followerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}
	user, err := GetUser(req.UserID)
	if err != nil {
		middleware.SendErrorResponse(w, "Error querying user data to struct"+err.Error(), http.StatusInternalServerError)
		return
	}
	status := "accepted"
	if user.ProfileVisibility == "private" {
		status = "pending"
	}

	if err := sqlite.Db.FollowUser(req.UserID, req.FollowerID, status); err != nil {
		middleware.SendErrorResponse(w, "Failed to save follower to DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	followerIDStr := r.URL.Query().Get("followerId")

	if userIDStr == "" || followerIDStr == "" {
		middleware.SendErrorResponse(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Convert userID and followerID to integers
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		middleware.SendErrorResponse(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	followerID, err := strconv.Atoi(followerIDStr)
	if err != nil {
		middleware.SendErrorResponse(w, "Invalid followerId parameter", http.StatusBadRequest)
		return
	}

	if err := sqlite.Db.UnfollowUser(userID, followerID); err != nil {
		middleware.SendErrorResponse(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func GetFollowStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	followerIDStr := r.URL.Query().Get("followerId")
	fmt.Println(userIDStr, followerIDStr)
	if userIDStr == "" || followerIDStr == "" {
		middleware.SendErrorResponse(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Convert userID and followerID to integers
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		middleware.SendErrorResponse(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	followerID, err := strconv.Atoi(followerIDStr)
	if err != nil {
		middleware.SendErrorResponse(w, "Invalid followerId parameter", http.StatusBadRequest)
		return
	}

	// Call the DB function to get the follow status
	isFollowing, isPending, err := sqlite.Db.CheckFollowStatus(userID, followerID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to retrieve follow status", http.StatusInternalServerError)
		return
	}

	// Respond with the follow status as JSON
	response := map[string]bool{
		"isFollowing": isFollowing,
		"isPending":   isPending,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPendingFollowRequests fetches pending follow requests for the logged-in user
func GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		middleware.SendErrorResponse(w, "Invalid userId", http.StatusBadRequest)
		return
	}
	pendingRequests, err := sqlite.Db.GetPendingFollowRequests(userID)
	if err != nil {
		middleware.SendErrorResponse(w, "Error fetching pending requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pendingRequests)
}

func AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID     int `json:"userId"`
		FollowerID int `json:"followerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}
	fmt.Println("checking", req.FollowerID, req.UserID)
	if err := sqlite.Db.UpdateFollowRequestStatus(req.UserID, req.FollowerID, "accepted"); err != nil {
		middleware.SendErrorResponse(w, "Failed to accept follow request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RejectFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID     int `json:"userId"`
		FollowerID int `json:"followerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := sqlite.Db.UpdateFollowRequestStatus(req.UserID, req.FollowerID, "rejected"); err != nil {
		middleware.SendErrorResponse(w, "Failed to reject follow request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func GetFollowLists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		middleware.SendErrorResponse(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Convert userID to integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		middleware.SendErrorResponse(w, "Invalid userId parameter", http.StatusBadRequest)
		return
	}

	// Get followers (users who follow the current user)
	followersIDs, err := sqlite.Db.GetFollowers(userID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to retrieve followers", http.StatusInternalServerError)
		return
	}

	// Get following (users that the current user follows)
	followingIDs, err := sqlite.Db.GetFollowing(userID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to retrieve following", http.StatusInternalServerError)
		return
	}

	// Fetch user info for followers and following users
	followersInfo, err := sqlite.Db.GetUsersByIDs(followersIDs)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to retrieve followers info", http.StatusInternalServerError)
		return
	}
	followingInfo, err := sqlite.Db.GetUsersByIDs(followingIDs)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to retrieve following info", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	responseData := map[string][]structs.User{
		"followers": followersInfo,
		"following": followingInfo,
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		middleware.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
