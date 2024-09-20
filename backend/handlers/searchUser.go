package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
)

func SearchUser(w http.ResponseWriter, r *http.Request) {
	// Parse the search query from the request (e.g., ?query=John)
	query := r.URL.Query().Get("query")
	if query == "" {
		middleware.SendErrorResponse(w, "Error getting token", http.StatusBadRequest)
		return
	}

	// Query the database to find matching users
	users, err := sqlite.Db.SearchUsersInDB(query)
	if err != nil {
		middleware.SendErrorResponse(w, "Error gettimg User(s) info from DB", http.StatusInternalServerError)
		return
	}

	// Return the users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
