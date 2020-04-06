package model

type Customer struct {
	Id       ID
	Name     string
	Cuit     int
	Projects ProjectList
}

type CustomerList []Customer
