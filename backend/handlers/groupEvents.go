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
	err := sqlite.Db.CreateEvent(req.GroupID, req.Title, req.Description, req.EventDate)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to insert event into DB"+err.Error(), http.StatusInternalServerError)
		return
	}
	//fire WS to send creation of new Group Event notification for group members online
	triggerGroupEventNotify(req)
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

func GetGroupMembersWithReactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	eventIDstr := r.URL.Query().Get("eventID")
	groupIDstr := r.URL.Query().Get("groupID")

	if eventIDstr == "" || groupIDstr == "" {
		middleware.SendErrorResponse(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	eventID, err := strconv.Atoi(eventIDstr)
	if err != nil || eventID <= 0 {
		middleware.SendErrorResponse(w, "Invalid eventID", http.StatusBadRequest)
		return
	}
	groupID, err := strconv.Atoi(groupIDstr)
	if err != nil || groupID <= 0 {
		middleware.SendErrorResponse(w, "Invalid groupID", http.StatusBadRequest)
		return
	}

	members, err := sqlite.Db.GetMembersWithReactions(eventID, groupID)
	if err != nil {
		middleware.SendErrorResponse(w, "Error fetching members with reactions", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(members)
}
