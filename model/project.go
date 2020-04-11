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

func (projectLis ProjectList) IDs() IDList {
	var ids IDList

	for _, project := range projectLis {
		ids = append(ids, project.ID)
	}

	return ids
}
