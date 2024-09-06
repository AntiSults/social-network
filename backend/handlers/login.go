package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/security"
	"social-network/structs"
	"sync"
)

// Creating local variable for storing users online.
var (
	UserMap     = make(map[int]structs.User)
	userMapLock sync.RWMutex // Mutex to protect UserMap
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := sqlite.Db.GetUserByEmail(email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				middleware.SendErrorResponse(w, "User email not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		// Compare passwords
		err = security.CheckPassword([]byte(user.Password), []byte(password))
		if err != nil {
			middleware.SendErrorResponse(w, "Wrong password", http.StatusBadRequest)
			return
		}

		security.NewSession("session_token", user.ID, w)

		// Protect UserMap with write lock
		userMapLock.Lock()
		UserMap[user.ID] = *user
		userMapLock.Unlock()

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
