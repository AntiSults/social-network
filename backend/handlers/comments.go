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
	"strconv"
	"time"
)

const commentsDir = "../public/postsContent"

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		middleware.SendErrorResponse(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		middleware.SendErrorResponse(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if content == "" {
		middleware.SendErrorResponse(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting token", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserId(cookie.Value)
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
		return
	}

	// Searching for comment creators name
	user, err := middleware.GetUser(userID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to retrieve user", http.StatusInternalServerError)
		return

	}

	comment := structs.Comment{
		PostID:          postID,
		Content:         content,
		CreatedAt:       time.Now(),
		UserID:          userID,
		AuthorFirstName: user.FirstName,
		AuthorLastName:  user.LastName,
	}

	if filePath, err := saveCommentFile(r); err == nil && filePath != "" {
		comment.File = filePath
	} else if err != nil && err != errNoFile {
		middleware.SendErrorResponse(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	if err := sqlite.Db.SaveComment(&comment); err != nil {
		middleware.SendErrorResponse(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		middleware.SendErrorResponse(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		middleware.SendErrorResponse(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	comments, err := sqlite.Db.GetComments(postID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func saveCommentFile(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	if err := os.MkdirAll(commentsDir, os.ModePerm); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
	filePath := filepath.Join(commentsDir, fileName)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", err
	}

	return filePath, nil
}
