package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/structs"
	"strconv"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		middleware.SendErrorResponse(w, "Invalid userID", http.StatusBadRequest)
		return
	}
	events, err := sqlite.Db.GetAllEvents(userID)
	if err != nil {
		middleware.SendErrorResponse(w, "Error fetching pending invitations", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
	var req structs.Event
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input"+err.Error(), http.StatusBadRequest)
		return
	}
	eventDate := "2024-09-30 12:00:00"
	eventDate = req.EventDate
	err := sqlite.Db.CreateEvent(req.GroupID, req.Title, req.Description, eventDate)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to insert group into DB"+err.Error(), http.StatusInternalServerError)
		return
	}
}
func EventReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		EventID  int    `json:"eventId"`
		UserID   int    `json:"userId"`
		Reaction string `json:"reaction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.SendErrorResponse(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := sqlite.Db.ReactToEvent(req.EventID, req.UserID, req.Reaction)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to process the invite request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Reaction processed successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
