package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
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
