package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
	"social-network/db/security"
	"social-network/db/sqlite"
	"social-network/middleware"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Get the hashed pw from db
		userID, hashedPw, err := sqlite.Db.GetId_Password(email)
		if err != nil {
			if err == sql.ErrNoRows {
				middleware.SendErrorResponse(w, "User email not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Compare passwords
		err = security.CheckPassword([]byte(hashedPw), []byte(password))
		if err != nil {
			middleware.SendErrorResponse(w, "Wrong password", http.StatusBadRequest)
			return
		}

		sessionToken, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Error creating session token", http.StatusInternalServerError)
			return
		}
		expiresAt := time.Now().Add(24 * time.Hour)

		err = sqlite.Db.SaveSession(userID, sessionToken, expiresAt)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error inserting a session", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken.String(),
			Expires:  expiresAt,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
