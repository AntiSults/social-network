package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
)

func SearchUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query().Get("query")
	if query == "" {
		middleware.SendErrorResponse(w, "Query parameter is required", http.StatusBadRequest)
		return
	}
	users, err := sqlite.Db.SearchUsersInDB(query)
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting User(s) info from DB", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID             int    `json:"userId"`
		ProfVisibility string `json:"profileVisibility"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err := sqlite.Db.UpdateProfileVisibility(req.ID, req.ProfVisibility)
	if err != nil {
		middleware.SendErrorResponse(w, "Error updating profile visibility", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{
		"message": "Profile visibility updated successfully",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		middleware.SendErrorResponse(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
