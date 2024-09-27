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

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
	var req structs.Group

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input"+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Description == "" || req.CreatorID == 0 {
		middleware.SendErrorResponse(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := sqlite.Db.CreateGroup(req.Name, req.Description, req.CreatorID)
	if err != nil {

		middleware.SendErrorResponse(w, "Failed to insert group into DB"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group created successfully"})
}

func GetGroupsWithMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}

	groups, err := sqlite.Db.GetGroupsWithMembers()
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func JoinGroupRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}

	var req struct {
		GroupID int `json:"groupId"`
		UserID  int `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err := sqlite.Db.RequestToJoinGroup(req.GroupID, req.UserID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to insert Joining request into DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Requested to join"})
}

func InviteToGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}

	var req struct {
		GroupID   int `json:"groupId"`
		UserID    int `json:"invitedUserId"`
		InviterID int `json:"inviterId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err := sqlite.Db.InviteUserToGroup(req.GroupID, req.UserID, req.InviterID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to insert Joining request into DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invite is sent"})
}

func JoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Define the expected request payload
	var req struct {
		GroupID int  `json:"groupId"`
		UserID  int  `json:"userId"`
		Accept  bool `json:"accept"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := sqlite.Db.HandleGroupRequest(req.GroupID, req.UserID, req.Accept)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to process the join request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	status := "rejected"
	if req.Accept {
		status = "accepted"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request processed successfully",
		"groupId": fmt.Sprintf("%d", req.GroupID),
		"userId":  fmt.Sprintf("%d", req.UserID),
		"status":  status,
	})
}

func InviteRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Define the expected request payload
	var req struct {
		GroupID int  `json:"groupId"`
		UserID  int  `json:"userId"`
		Accept  bool `json:"accept"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Invites", req)
	err := sqlite.Db.HandleGroupRequest(req.GroupID, req.UserID, req.Accept)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to process the invite request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	status := "rejected"
	if req.Accept {
		status = "accepted"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request processed successfully",
		"groupId": fmt.Sprintf("%d", req.GroupID),
		"userId":  fmt.Sprintf("%d", req.UserID),
		"status":  status,
	})
}

// GetPendingGroupJoin fetches pending group join requests for the group creator
func GetPendingGroupJoin(w http.ResponseWriter, r *http.Request) {
	creatorStr := r.URL.Query().Get("creatorID")
	creatorID, err := strconv.Atoi(creatorStr)
	if err != nil || creatorID <= 0 {
		middleware.SendErrorResponse(w, "Invalid creatorID", http.StatusBadRequest)
		return
	}
	pendingRequests, err := sqlite.Db.GetPendingGroupRequests(creatorID)
	if err != nil {
		middleware.SendErrorResponse(w, "Error fetching pending requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pendingRequests)
}

func GetPendingGroupInvites(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		middleware.SendErrorResponse(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	invitations, err := sqlite.Db.GetPendingGroupInvites(userID)
	if err != nil {
		middleware.SendErrorResponse(w, "Error fetching pending invitations", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(invitations)
}
