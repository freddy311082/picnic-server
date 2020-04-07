package model

type User struct {
	Id       ID
	Name     string
	LastName string
	Email    string
	Token    string
}

type UserList []*User
