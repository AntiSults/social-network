package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social-network/backend/db/sqlite"
	"social-network/backend/middleware"
	"social-network/backend/structs"
)

func GetUserData(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {

		cookie, err := r.Cookie("session_token")
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting token" + err.Error(), http.StatusBadRequest)
			return
		}

		db, err := sqlite.OpenDatabase()
		if err != nil {
			middleware.SendErrorResponse(w, "Error opening database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		userID, err := sqlite.GetUserIdFromToken(db, cookie.Value)
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
			return
		}

		var user structs.User
		var nick sql.NullString
		var aboutMe sql.NullString
		var avatarPath sql.NullString
		err = db.QueryRow("SELECT Email, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath FROM Users WHERE ID = ?", userID).Scan(&user.Email, &user.FirstName, &user.LastName, &user.DOB, &nick, &aboutMe, &avatarPath)
		if err != nil {
			middleware.SendErrorResponse(w, "Error querying user data to struct"+err.Error(), http.StatusInternalServerError)
			return
		}

		if nick.Valid {
			user.NickName = nick.String
		}
		if aboutMe.Valid {
			user.AboutMe = aboutMe.String
		}
		if avatarPath.Valid {
			user.AvatarPath = avatarPath.String
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			middleware.SendErrorResponse(w, "Error marshalling user data to JSON", http.StatusInternalServerError)
		return

		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) 
		w.Write(jsonData)
	} else {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}