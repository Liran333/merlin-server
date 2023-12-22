package controller

type User struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	AvatarId    string `json:"avatar_id"`
	Email       string `json:"email"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	Website     string `json:"website"`
	Owner       string `json:"owner"`
}
