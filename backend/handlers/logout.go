package handlers

import (
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/security"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
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
			userMapLock.Lock()
			delete(UserMap, session.UserID)
			userMapLock.Unlock()
		}

		security.RemoveSession(cookie.Value)

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
