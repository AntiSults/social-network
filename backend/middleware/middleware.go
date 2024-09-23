package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"social-network/db/sqlite"
	"social-network/structs"
)

// Allows CORS from specific origin
func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Headers for CORS
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
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
		userID, err := sqlite.Db.GetUserIdFromToken(cookie.Value)
		if err != nil {
			fmt.Println(err)
			SendErrorResponse(w, "Please log in. "+err.Error(), http.StatusUnauthorized)
			return
		}
		fmt.Println(userID)
		//?
		next.ServeHTTP(w, r)
	})
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	redirectURL := "/login?error=" + url.QueryEscape(message)
	http.Redirect(w, r, redirectURL, statusCode)
}

func SendErrorResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(structs.ErrorResponse{Message: message})
}

// TEST FUNCTION DELETE LATER
func DummyCheck(w http.ResponseWriter, r *http.Request) {

	fmt.Println("joujou")
	w.WriteHeader(http.StatusOK)
}
