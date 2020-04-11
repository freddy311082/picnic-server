package api

import (
	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/service"
	"github.com/graphql-go/graphql"
	"time"
)

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutations",
	Fields: graphql.Fields{
		"registerUser": &graphql.Field{
			Type: UserType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "User Name",
				},
				"lastName": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "User Last Name",
				},
				"email": &graphql.ArgumentConfig{
					Type:        &graphql.NonNull{OfType: graphql.String},
					Description: "User email: This field is mandatory as it will be the username used to login.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				user, err := service.Instance().RegisterUser(&model.User{
					Name:     p.Args["name"].(string),
					LastName: p.Args["lastName"].(string),
					Email:    p.Args["email"].(string),
				})

				if err != nil {
					return nil, err
				}

				return gqlUserFromModel(user), err
			},
			Description: "Register a new user in the system by email. If the user already exists and error will" +
				" be raised.",
		},
		"createProject": &graphql.Field{
			Type: ProjectType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        &graphql.NonNull{OfType: graphql.String},
					Description: "Project name. It cannot be null.",
				},
				"description": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Description about the projects.",
				},
				"created_at": &graphql.ArgumentConfig{
					Type:         &graphql.NonNull{OfType: graphql.DateTime},
					DefaultValue: time.Now(),
					Description:  "",
				},
			},
		},
	},
	Description: "Mutations definitions for Picnic GraphQL API.",
})
