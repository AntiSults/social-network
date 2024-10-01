package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/db/sqlite"
	"social-network/security"
	"social-network/structs"
	"sync"
)

// Creating local variable for storing users online.
var (
	UserMap     = make(map[int]structs.User)
	UserMapLock sync.RWMutex // Mutex to protect UserMap
)

// Allows CORS from specific origin
func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Headers for CORS
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, DELETE")
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
			// fmt.Println(err)
			SendErrorResponse(w, "Invalid session. Please log in.", http.StatusUnauthorized)
			return
		}
		// Validate the session
		if !security.ValidateSession(cookie.Value) {
			SendErrorResponse(w, "Please log in.", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SendErrorResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(structs.ErrorResponse{Message: message})
}

// GetUserId is getting user ID with token, either from sessions map or from DB
func GetUserId(token string) (int, error) {
	userID := 0
	var err error

	// Acquire a read lock before accessing the shared map
	security.SessionLock.RLock()
	session, ok := security.DbSessions[token]
	security.SessionLock.RUnlock()

	if ok {
		userID = session.UserID
	} else {
		// Fall back to database lookup if not found in in-memory store
		userID, err = sqlite.Db.GetUserIdFromToken(token)
		if err != nil {
			return -1, fmt.Errorf("error getting ID from session token: %w", err)
		}
	}
	return userID, err
}

// GetUser is getting user with user ID, either from User map or from DB
func GetUser(id int) (*structs.User, error) {
	var (
		user *structs.User
		err  error
	)
	UserMapLock.RLock()
	u, ok := UserMap[id]
	UserMapLock.RUnlock()

	if ok {
		user = &u
	} else {
		// Fall back to database lookup if not found in in-memory store
		user, err = sqlite.Db.GetUser(id)
		if err != nil {
			return nil, fmt.Errorf("error querying user data to struct:  %w", err)
		}
	}
	return user, nil
}

// func RedirectToLogin(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
// 	redirectURL := "/login?error=" + url.QueryEscape(message)
// 	http.Redirect(w, r, redirectURL, statusCode)
// }

// // TEST FUNCTION DELETE LATER
// func DummyCheck(w http.ResponseWriter, r *http.Request) {

// 	fmt.Println("joujou")
// 	w.WriteHeader(http.StatusOK)
// }
