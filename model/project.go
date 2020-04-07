package model

import "time"

type Project struct {
	ID          ID
	Name        string
	Description string
	CreatedAt   time.Time
	Owner       *User
	Customer    *Customer
	Fields      ProjectFieldList
}

type ProjectList []*Project

type ProjectField struct {
	Name string
}

type ProjectFieldList []ProjectField
