package api

import (
	"github.com/freddy311082/picnic-server/model"
	"github.com/graphql-go/graphql"
)

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},

		"name": &graphql.Field{
			Type: graphql.String,
		},

		"lastName": &graphql.Field{
			Type: graphql.String,
		},

		"email": &graphql.Field{
			Type: graphql.String,
		},
	},
	Description: "User object type definition.",
})

type gqlUserResponse struct {
	ID       string
	Name     string
	LastName string
	Email    string
}

type gqlUserListResponse []*gqlUserResponse

func gqlUserFromModel(user *model.User) *gqlUserResponse {
	return &gqlUserResponse{
		ID:       user.Id.ToString(),
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
	}
}

func gqlUserListFromModel(userList model.UserList) gqlUserListResponse {
	var gqlUserList gqlUserListResponse

	for _, user := range userList {
		gqlUserList = append(gqlUserList, gqlUserFromModel(&user))
	}

	return gqlUserList
}
