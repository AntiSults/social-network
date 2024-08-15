package structs

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	DOB string `json:"dob"`
	NickName string `json:"nickname"`
	AboutMe string `json:"aboutMe"`
	AvatarPath string `json:"avatarPath"`
	ProfileVisibility string `json:"profileVisibility"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}