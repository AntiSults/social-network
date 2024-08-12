package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"social-network/backend/db/sqlite"
	"social-network/backend/middleware"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		
		email:= r.FormValue("email")
		password := r.FormValue("password")

		db, err := sqlite.OpenDatabase()
		if err != nil {
			http.Error(w, "Couldn't connect to database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Get the hashed pw from db
		var userID string
		var hashedPw string
		err = db.QueryRow("SELECT ID, Password FROM Users WHERE Email = ?", email).Scan(&userID, &hashedPw)
		if err != nil {
			if err == sql.ErrNoRows {
				middleware.SendErrorResponse(w, "User email not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Compare passwords 
		err = bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(password))
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

		_, err = db.Exec(`
			INSERT INTO Sessions (UserID, SessionToken, ExpiresAt) VALUES (?, ?, ?)
		`, userID, sessionToken, expiresAt)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error inserting a session", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name: "session_token",
			Value: sessionToken.String(),
			Expires: expiresAt,
			SameSite: http.SameSiteNoneMode,
			Secure: true,
		}
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}