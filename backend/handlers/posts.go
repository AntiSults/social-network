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

const postsDir = "../public/postsContent"

func CreatePost(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		middleware.SendErrorResponse(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20) // 10MB limit

	content := r.FormValue("content")
	privacy := r.FormValue("privacy")
	groupIDStr := r.FormValue("group_id")
	visibleUsersStr := r.FormValue("visible_users")

	var post structs.Post
	post.Content = content
	post.Privacy = privacy

	cookie, err := r.Cookie("session_token")
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting token"+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserId(cookie.Value)
	if err != nil {
		middleware.SendErrorResponse(w, "Error getting ID from session token", http.StatusInternalServerError)
		return
	}

	var groupID int
	if groupIDStr != "" {
		parsedGroupID, err := strconv.Atoi(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		groupID, err = sqlite.Db.GetGroupIDForUser(userID, parsedGroupID)
		if err != nil {
			http.Error(w, "Error checking group membership", http.StatusInternalServerError)
			return
		}

		if groupID == 0 {
			http.Error(w, "User is not a member of the specified group", http.StatusForbidden)
			return
		}
	}

	// Searching for post creators name
	user, err := middleware.GetUser(userID)
	if err != nil {
		middleware.SendErrorResponse(w, "Failed to retrieve user", http.StatusInternalServerError)
		return

	}

	post.UserID = userID
	post.CreatedAt = time.Now()
	post.AuthorFirstName = user.FirstName
	post.AuthorLastName = user.LastName
	post.GroupID = groupID

	if post.Privacy == "" {
		post.Privacy = "public"
	}

	if post.Privacy == "almost private" {
		if visibleUsersStr != "" {
			var visibleUsers []int
			err = json.Unmarshal([]byte(visibleUsersStr), &visibleUsers)
			if err != nil {
				middleware.SendErrorResponse(w, "Failed to parse visible users: "+err.Error(), http.StatusBadRequest)
				return
			}
			post.VisibleUsers = visibleUsers
		} else {
			post.VisibleUsers = []int{}
		}
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
		posts, err = sqlite.Db.GetPosts(0, false)
	} else {
		userID, err := middleware.GetUserId(cookie.Value)
		if err != nil {
			posts, err = sqlite.Db.GetPosts(0, false)
		} else {
			posts, err = sqlite.Db.GetPosts(userID, true)
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
	file, handler, err := r.FormFile("files")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
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
}
