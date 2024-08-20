package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"social-network/db/sqlite"
	"social-network/security"
	"social-network/structs"
	"time"
)

type RegisterForm struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Dob       string `json:"dob"`
}

const avatarDir = "../public/uploads"

var errNoFile = fmt.Errorf("no file")

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user := structs.User{
			Email:     r.FormValue("email"),
			FirstName: r.FormValue("firstName"),
			LastName:  r.FormValue("lastName"),
			NickName:  r.FormValue("nickname"),
			AboutMe:   r.FormValue("aboutMe"),
		}

		dob := r.FormValue("dob")

		parsedDob, err := time.Parse("2006-01-02", dob)
		if err != nil {
			http.Error(w, "Failed to parse date", http.StatusBadRequest)
			return
		}
		user.DOB = parsedDob

		hashedPw, err := security.HashPassword(r.FormValue("password"))
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPw)

		avatarPath, err := saveImage(r)
		if err != nil {
			if !errors.Is(err, errNoFile) {
				http.Error(w, "Failed to save image file", http.StatusInternalServerError)
				return
			}
		} else {
			user.AvatarPath = avatarPath
		}

		err = sqlite.Db.InsertUserToDatabase(user)
		if err != nil {
			http.Error(w, "Failed to insert into Users table", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func saveImage(r *http.Request) (string, error) {
	if file, handler, err := r.FormFile("avatar"); err == nil {

		defer file.Close()

		// Makes sure the directory exists
		if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
			return "", err
		}

		// Creates a specific file name for the img
		avatarFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		avatarPath := filepath.Join(avatarDir, avatarFileName)

		// Creates the file in the server
		img, err := os.Create(avatarPath)
		if err != nil {
			return "", err
		}
		defer img.Close()

		// Copies the file to the file server
		if _, err := io.Copy(img, file); err != nil {
			return "", err
		}
		return avatarPath, nil
	} else {
		return "", errNoFile
	}

}
