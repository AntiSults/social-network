package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
)

func GetAvatarPath(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	logged := false
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err != http.ErrNoCookie {
			middleware.SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		logged = true
	}
	if logged {
		userID, err := middleware.GetUserId(cookie.Value)
		if err != nil {
			middleware.SendErrorResponse(w, "error getting id from token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		avatarPath, err := sqlite.Db.GetAvatarFromID(userID)
		if err != nil {
			avatarPath = ""
		}
		response := map[string]string{
			"avatarPath": avatarPath,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
