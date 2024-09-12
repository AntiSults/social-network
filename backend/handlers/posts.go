package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/structs"
	"time"
)

const postsDir = "../public/postsContent"

func CreatePost(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20) // 10MB limit

		content := r.FormValue("content")
		privacy := r.FormValue("privacy")

		fmt.Println("Content:", content)
		fmt.Println("Privacy:", privacy)

		var post structs.Post
		post.Content = content
		post.Privacy = privacy

		cookie, err := r.Cookie("session_token")
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting token"+err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := GetUserId(cookie.Value)
		if err != nil {
			middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
			return
		}

		// Searching for post creators name
		user, err := GetUser(userID)
		if err != nil {
			middleware.SendErrorResponse(w, "Failed to retrieve user", http.StatusInternalServerError)
			return

		}

		post.UserID = userID
		post.CreatedAt = time.Now()
		post.AuthorFirstName = user.FirstName
		post.AuthorLastName = user.LastName

		if post.Privacy == "" {
			post.Privacy = "public"
		}

		filePath, err := savePostFile(r)
		if err != nil && err != errNoFile {
			middleware.SendErrorResponse(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if filePath != "" {
			post.Files = filePath
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

		// Attempt to retrieve the session cookie
		cookie, err := r.Cookie("session_token")

		if err != nil || cookie.Value == "" {
			// If there is no valid cookie, show only public posts
			posts, err = sqlite.Db.GetPosts(false)
		} else {
			// If a cookie exists, check user authentication
			_, err = GetUserId(cookie.Value)
			if err != nil {
				// If user is not authenticated, show only public posts
				posts, err = sqlite.Db.GetPosts(false)
			} else {
				// If user is authenticated, show all posts
				posts, err = sqlite.Db.GetPosts(true)
			}
		}

		// Handle potential error from database query
		if err != nil {
			middleware.SendErrorResponse(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Set response header and encode posts as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)

	} else {
		// Method not allowed error for non-GET requests
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}
}

func savePostFile(r *http.Request) (string, error) {
	if file, handler, err := r.FormFile("files"); err == nil {
		defer file.Close()

		if err := os.MkdirAll(postsDir, os.ModePerm); err != nil {
			return "", err
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		filePath := filepath.Join(postsDir, fileName)

		outFile, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, file); err != nil {
			return "", err
		}

		return filePath, nil
	} else {
		return "", err
	}
}
