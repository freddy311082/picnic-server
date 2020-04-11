package model

type User struct {
	ID       ID
	Name     string
	LastName string
	Email    string
	Token    string
}

type UserList []*User
