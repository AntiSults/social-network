package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterForm struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Dob       string `json:"dob"`
}

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
		var form RegisterForm
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		fmt.Println(form)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Registration successful"))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
