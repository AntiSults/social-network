package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"social-network/backend/db/sqlite"
	"social-network/backend/middleware"
	"social-network/backend/structs"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

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
			ProfileVisibility: "public",
    	}
		emailRegex := `^[\w\.-]+@[\w\.-]+\.[a-zA-Z]{2,}$`
		match, err := regexp.Match(emailRegex, []byte(user.Email))
		if err != nil {
			middleware.SendErrorResponse(w, "error with regex", http.StatusBadRequest)
			return
		}
		if (!match) {
			middleware.SendErrorResponse(w, "Invalid email aadress", http.StatusBadRequest)
			return
		}

		dob := r.FormValue("dob")

		parsedDob, err := time.Parse("2006-01-02", dob)
		if err != nil {
			middleware.SendErrorResponse(w, "Failed to parse date", http.StatusBadRequest)
			return
		}
		user.DOB = parsedDob.Format("2006-01-02")

		password := r.FormValue("password")

		if(len(password) < 4) {
			middleware.SendErrorResponse(w, "Password too short!", http.StatusBadRequest)
		}

		hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			middleware.SendErrorResponse(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPw)

		avatarPath, err := saveImage(r)
		if err != nil {
			if !errors.Is(err, errNoFile) {
				middleware.SendErrorResponse(w, "Failed to save image file", http.StatusInternalServerError)
				return
			} 
		} else {
			user.AvatarPath = avatarPath
		}

		err = insertUserToDatabase(user)
		if err !=nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					middleware.SendErrorResponse(w, "Email already exists", http.StatusBadRequest)
					return
				}
			}
			middleware.SendErrorResponse(w, "Failed to insert into Users table" + err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		middleware.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func saveImage (r *http.Request) (string, error) {
	if file, handler, err := r.FormFile("avatar"); err==nil {
		
		defer file.Close()

		// Makes sure the directory exists
		if err := os.MkdirAll(avatarDir, os.ModePerm); err!= nil {
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
		if _, err := io.Copy(img, file); err!=nil{
			return "", err
		}
		return avatarPath, nil
	} else {
		return "", errNoFile
	}
	
}

func insertUserToDatabase(user structs.User) error {

	db, err := sqlite.OpenDatabase()
	if err !=nil  {return err}
	defer db.Close()

	prep, err := db.Prepare(`
			INSERT INTO Users (Email, Password, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`)
	if err !=nil  {
		return err
	}
	defer prep.Close()

	var nickName, aboutMe, avatarPath interface{}
	if user.NickName == "" {
		nickName = nil
	} else {
		nickName = user.NickName
	}

	if user.AboutMe == "" {
		aboutMe = nil
	} else {
		aboutMe = user.AboutMe
	}

	if user.AvatarPath == "" {
		avatarPath = nil
	} else {
		avatarPath = user.AvatarPath
	}

	_, err = prep.Exec(user.Email, user.Password, user.FirstName, user.LastName, user.DOB, nickName, aboutMe, avatarPath)
	if err !=nil  {
		return err
	}
	fmt.Println("Successfully inserted to db!")

	return nil
}
