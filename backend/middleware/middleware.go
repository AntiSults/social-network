package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"social-network/backend/db/sqlite"
	"social-network/backend/structs"
)

// Allows CORS from specific origin
func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Headers for CORS
		if origin == "http://localhost:3000" { 
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight (OPTIONS) requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
        if err != nil {
			fmt.Println(err)
			SendErrorResponse(w, "Invalid session. Please log in.", http.StatusUnauthorized)
			return
        }

        // Validate the session
        userID, check, err := validateSession(cookie.Value)
        if !check {
			fmt.Println(err)
            SendErrorResponse(w, "Please log in. " + err.Error(), http.StatusUnauthorized)
            return
        }
		fmt.Println(userID)
//?
		next.ServeHTTP(w, r)
	})
}

// Compares current session from cookies against database
func validateSession(sessionID string) (int, bool, error) {
	db, err := sqlite.OpenDatabase()
	if err != nil {
		return -1, false, err
	}
	defer db.Close()
	var userID int
	err = db.QueryRow(`
		SELECT UserID FROM Sessions WHERE SessionToken = ? 
	`, sessionID).Scan(&userID)
	if err != nil {
		return -1, false, err
	}
	return userID, true, nil
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	redirectURL := "/login?error=" + url.QueryEscape(message)
	http.Redirect(w, r, redirectURL, statusCode)
}

func SendErrorResponse(w http.ResponseWriter, message string, code int){
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(structs.ErrorResponse{Message: message})
}

// TEST FUNCTION DELETE LATER
func DummyCheck(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("joujou")
	w.WriteHeader(http.StatusOK)
}