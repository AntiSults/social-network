package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"social-network/backend/db/sqlite"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*") 
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
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
				http.Error(w, "User email not found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Compare passwords 
		err = bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(password))
		if err != nil {
			http.Error(w, "Wrong password", http.StatusUnauthorized)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}