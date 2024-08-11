package middleware

import (
	"fmt"
	"net/http"
	"social-network/backend/db/sqlite"
)

func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Allow only specific origins 
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

/*NEED TO SEND ERROR NOT REDIRECT 
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid session. Please log in again."})
			return
NEED TO HANDLE ERRORS ON FRONTEND LOGINS/PAGE.TSX
	else {
        // Handle errors by reading the error message from the response
        const data = await response.json();
        setError(data.message || 'Login failed');
      }			 
*/

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
        if err != nil {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        // Validate the session
        userID, check, err := validateSession(cookie.Value)
        if !check {
			fmt.Println(err)
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }
		fmt.Println(userID)
	})
}

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

func DummyCheck(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("joujou")
	w.WriteHeader(http.StatusOK)
}