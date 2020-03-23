package api

import "github.com/graphql-go/graphql"

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},

		"name": &graphql.Field{
			Type: graphql.String,
		},

		"last_name": &graphql.Field{
			Type: graphql.String,
		},

		"email": &graphql.Field{
			Type: graphql.String,
		},
	},
	Description: "User object type definition.",
})
