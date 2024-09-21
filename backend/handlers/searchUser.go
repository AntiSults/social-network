package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
)

func SearchUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		middleware.SendErrorResponse(w, "Query parameter is required", http.StatusBadRequest)
		return
	}
	// Query the database to find matching users
	users, err := sqlite.Db.SearchUsersInDB(query)
	if err != nil {
		middleware.SendErrorResponse(w, "Error gettimg User(s) info from DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
