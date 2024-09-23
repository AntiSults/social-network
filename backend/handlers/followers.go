package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
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
