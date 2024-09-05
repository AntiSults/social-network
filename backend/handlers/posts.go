package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/structs"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var post structs.Post

		cookie, err := r.Cookie("session_token")
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting token"+err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {

			middleware.SendErrorResponse(w, "Invalid input", http.StatusInternalServerError)
			return
		}

		userID, err := sqlite.Db.GetUserIdFromToken(cookie.Value)
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
			return
		}

		// Searching for post creators name
		firstName, lastName, err := sqlite.Db.GetUserNameByID(userID)
		if err != nil {

			middleware.SendErrorResponse(w, "Failed to retrieve user details", http.StatusInternalServerError)
			return

		}

		post.UserID = userID
		post.CreatedAt = time.Now()
		post.AuthorFirstName = firstName
		post.AuthorLastName = lastName

		if post.Privacy == "" {
			post.Privacy = "public"
		}

		if err := sqlite.Db.SavePost(&post); err != nil {
			middleware.SendErrorResponse(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	} else {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		var posts []structs.Post
		var err error
		cookie, err := r.Cookie("session_token")
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting token"+err.Error(), http.StatusBadRequest)
			return
		}
		_, err = sqlite.Db.GetUserIdFromToken(cookie.Value)

		if err != nil {
			// If user is not authenticated then show only public posts
			posts, err = sqlite.Db.GetPosts(false)
		} else {
			// If user is authenticated then show all posts
			posts, err = sqlite.Db.GetPosts(true)
		}

		if err != nil {
			middleware.SendErrorResponse(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)

	} else {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}
