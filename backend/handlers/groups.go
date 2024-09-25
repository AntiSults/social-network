package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/structs"
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
	fmt.Println("Checking", req)
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
func JoinRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		GroupID int  `json:"groupId"`
		UserID  int  `json:"userId"`
		Accept  bool `json:"accept"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := sqlite.Db.HandleGroupRequest(req.GroupID, req.UserID, req.Accept)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to insert respond to Join Request into DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reacted fo join request"})

}
