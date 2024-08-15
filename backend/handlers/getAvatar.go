package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social-network/backend/db/sqlite"
	"social-network/backend/middleware"
)

func GetAvatarPath(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
			userID, err := sqlite.Db.GetUserIdFromToken(cookie.Value)
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
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getAvatarFromID(db *sql.DB, id int) (string, error) {
	var avatarPath string
	err := db.QueryRow("SELECT AvatarPath FROM Users WHERE ID = ?", id).Scan(&avatarPath)
	if err != nil {
		return "", err
	}
	return avatarPath, nil
}
