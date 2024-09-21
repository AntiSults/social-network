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

	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
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

}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
	}

	var posts []structs.Post
	var err error

	cookie, err := r.Cookie("session_token")

	if err != nil || cookie.Value == "" {
		posts, err = sqlite.Db.GetPosts(false)
	} else {
		_, err = GetUserId(cookie.Value)
		if err != nil {
			posts, err = sqlite.Db.GetPosts(false)
		} else {
			posts, err = sqlite.Db.GetPosts(true)
		}
	}

	if err != nil {
		middleware.SendErrorResponse(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
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
