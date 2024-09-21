package handlers

import (
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/security"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting token: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = sqlite.Db.DeleteSessionFromDB(cookie.Value)
	if err != nil {
		if err.Error() != "no rows" {
			middleware.SendErrorResponse(w, "Error deleting from database: "+err.Error(), http.StatusBadRequest)
			return
		}
		middleware.SendErrorResponse(w, "No rows to delete: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if session exists in memory
	if session, ok := security.DbSessions[cookie.Value]; ok {
		// Protect UserMap with write lock before deletion
		UserMapLock.Lock()
		delete(UserMap, session.UserID)
		UserMapLock.Unlock()
	}

	security.RemoveSession(cookie.Value)

	// Delete cookie from the client side
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   "",
		Expires: time.Unix(0, 0), // Set expiry to a time in the past
		MaxAge:  -1,              // Also use MaxAge=-1 to ensure deletion
	})
	w.WriteHeader(http.StatusOK)

}
