package api

import "github.com/graphql-go/graphql"

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"users": &graphql.Field{
			Type: &graphql.List{
				OfType: UserType,
			},
			Args: graphql.FieldConfigArgument{
				"start_pos": &graphql.ArgumentConfig{
					Type:         &graphql.NonNull{OfType: graphql.Int},
					DefaultValue: 0,
					Description:  "",
				},
				"offset": &graphql.ArgumentConfig{
					Type:         &graphql.NonNull{OfType: graphql.Int},
					DefaultValue: 0,
					Description: `Number of users per page. By default, the number is 0. If 0 is passed, then all 
users will be returned. If a positive number is passed, then the amount of users returned will be less or equal than
the offset.`,
				},
			},
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {

			},
		},
	},
	Description: "Root Query for Picnic GraphQL Web Server",
})