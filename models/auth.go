package models

type UserRegData struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UserLoginData struct {
	Email    string
	Password string
}

type UserInfo struct {
	FirstName string `json:"first_name,string"`
	LastName  string `json:"last_name,string"`
	Email     string `json:"email,string"`
}
