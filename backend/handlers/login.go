package handlers

import (
	"database/sql"
	"net/http"
	"social-network/backend/db/sqlite"

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
		var hashedPw string
		err = db.QueryRow("SELECT Password FROM Users WHERE Email = ?", email).Scan(&hashedPw)
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
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}