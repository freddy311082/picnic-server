package model

type ID interface {
	ToString() string
}

type IDList []ID
