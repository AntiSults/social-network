package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust for production
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		// Handle preflight requests
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("Made it here")
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
    	password := r.FormValue("password")
    	firstName := r.FormValue("firstName")
    	lastName := r.FormValue("lastName")
    	dob := r.FormValue("dob")
		fmt.Println(email, password, firstName, lastName, dob)
		if file, handler, err := r.FormFile("avatar"); err==nil {
			defer file.Close()

			if err := os.MkdirAll(avatarDir, os.ModePerm); err!= nil {
				http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
            	return
			}
			fmt.Println("Handler: " + handler.Filename)
			avatarFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
        	avatarPath := filepath.Join(avatarDir, avatarFileName)

			img, err := os.Create(avatarPath)
			if err != nil {
				http.Error(w, "Failed to save the file", http.StatusInternalServerError)
            	return
			}
			defer img.Close()

			if _, err := io.Copy(img, file); err!=nil{
				http.Error(w, "Failed to save the file", http.StatusInternalServerError)
            	return
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
