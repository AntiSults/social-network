package structs

import "time"

type User struct {
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	DOB        time.Time `json:"dob"`
	NickName   string    `json:"nickname"`
	AboutMe    string    `json:"aboutMe"`
	AvatarPath string    `json:"avatarPath"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Message struct {
	Content string `json:"content"`
	Created string `json:"created"`
}
