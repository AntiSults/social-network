package handlers

import (
	"database/sql"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/security"
	"social-network/structs"
)

var User = structs.User{}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := sqlite.Db.GetUser_By_email(email)
		if err != nil {
			if err == sql.ErrNoRows {
				middleware.SendErrorResponse(w, "User email not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		User = *user
		// Compare passwords
		err = security.CheckPassword([]byte(User.Password), []byte(password))
		if err != nil {
			middleware.SendErrorResponse(w, "Wrong password", http.StatusBadRequest)
			return
		}

		security.NewSession("session_token", user.ID, w)

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
