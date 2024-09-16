package structs

import (
	"time"
)

type User struct {
	ID                int    `json:"ID,omitempty"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	DOB               string `json:"dob"`
	NickName          string `json:"nickname"`
	AboutMe           string `json:"aboutMe"`
	AvatarPath        string `json:"avatarPath"`
	ProfileVisibility string `json:"profileVisibility"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Message struct {
	ID          int    `json:"id,omitempty"`
	Content     string `json:"content"`
	Created     string `json:"created"`
	RecipientID []int  `json:"toUserID"`
	SenderID    int    `json:"fromUserID"`
}

type ChatMessage struct {
	Message []Message
	User    []User
}

type Post struct {
	ID              int       `json:"id,omitempty"`
	UserID          int       `json:"user_id"`
	Content         string    `json:"content"`
	Privacy         string    `json:"privacy"`
	CreatedAt       time.Time `json:"created_at"`
	AuthorFirstName string    `json:"author_first_name"`
	AuthorLastName  string    `json:"author_last_name"`
	Files           string    `json:"files"`
}
