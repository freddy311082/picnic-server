package model

type Customer struct {
	ID       ID
	Name     string
	Cuit     string
	Projects ProjectList
}

type CustomerList []*Customer
