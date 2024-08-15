package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"social-network/backend/db/sqlite"
	"social-network/backend/middleware"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting token"+err.Error(), http.StatusBadRequest)
			return
		}

		err = sqlite.Db.DeleteSessionFromDB(cookie.Value)
		if err != nil {
			if err.Error() != "no rows" {
				middleware.SendErrorResponse(w, "Error deleting from database"+err.Error(), http.StatusBadRequest)
				return
			}
			middleware.SendErrorResponse(w, "No rows to delete"+err.Error(), http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func deleteSessionFromDB(db *sql.DB, session string) error {
	stmt, err := db.Prepare("DELETE FROM Sessions WHERE SessionToken = ?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(session)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows")
	}
	return nil
}
