package model

type ID interface {
	ToString() string
}

type User struct {
	Id       ID
	Name     string
	LastName string
	Email    string
	Token    string
}

type UserList []User
