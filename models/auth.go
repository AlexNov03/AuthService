package models

type UserRegData struct {
	FirstName string `json:"firstname" validate:"required,min=1"`
	LastName  string `json:"lastname" validate:"required,min=1"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type UserLoginData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserInfo struct {
	ID        string `json:"user_id,string"`
	FirstName string `json:"first_name,string"`
	LastName  string `json:"last_name,string"`
	Email     string `json:"email,string"`
}
