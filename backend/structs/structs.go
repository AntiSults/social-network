package structs

import "time"

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	DOB time.Time `json:"dob"`
	NickName string `json:"nickname"`
	AboutMe string `json:"aboutMe"`
	AvatarPath string `json:"avatarPath"`
}