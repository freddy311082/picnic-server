package model

type User struct {
	ID       string
	Name     string
	LastName string
	Email    string
	Token    string
}

type UserList []User
