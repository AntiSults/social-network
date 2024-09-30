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

type Group struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorID   int       `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
	Members     []int     `json:"members"`
}

type Event struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"eventDate"`
	GroupID     int    `json:"groupId"`
}

type GroupJoinRequest struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	UserID    int    `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Message struct {
	ID          int    `json:"id,omitempty"`
	Content     string `json:"content"`
	Created     string `json:"created"`
	RecipientID int    `json:"toUserID"`
	SenderID    int    `json:"fromUserID"`
	GroupID     int    `json:"groupID"`
}

// type GroupMessage struct {
// 	ID          int    `json:"id,omitempty"`
// 	Content     string `json:"content"`
// 	Created     string `json:"created"`
// 	RecipientID int    `json:"toUserID"`
// 	SenderID    int    `json:"fromUserID"`
// 	GroupID     int    `json:"groupID"`
// }

type ChatMessage struct {
	Message []Message
	User    []User
	// Group   []Group
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
