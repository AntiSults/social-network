package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"social-network/db/sqlite"
	"social-network/structs"
	"time"
)

func CreatePost(db *sqlite.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post structs.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Getting userID from session
		userID, err := GetUserIDFromSession(r, db)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Searching for post creators name
		firstName, lastName, err := db.GetUserNameByID(userID)
		if err != nil {
			http.Error(w, "Failed to retrieve user details", http.StatusInternalServerError)
			return
		}

		post.UserID = userID
		post.CreatedAt = time.Now()
		post.AuthorFirstName = firstName
		post.AuthorLastName = lastName

		if post.Privacy == "" {
			post.Privacy = "public"
		}

		if err := db.SavePost(&post); err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

func GetUserIDFromSession(r *http.Request, db *sqlite.Database) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, errors.New("error getting session token: " + err.Error())
	}

	userID, err := db.GetUserIdFromToken(cookie.Value)
	if err != nil {
		return 0, errors.New("error getting user ID from session token: " + err.Error())
	}

	return userID, nil
}

func GetPosts(db *sqlite.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var posts []structs.Post
		var err error

		_, err = GetUserIDFromSession(r, db)
		if err != nil {
			// If user is not authenticated then show only public posts
			posts, err = db.GetPosts(false)
		} else {
			// If user is authenticated then show all posts
			posts, err = db.GetPosts(true)
		}

		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
