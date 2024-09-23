package handlers

import (
	"encoding/json"
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
