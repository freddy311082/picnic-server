package model

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type UserList []User
