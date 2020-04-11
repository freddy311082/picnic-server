package api

import (
	"github.com/freddy311082/picnic-server/service"
	"github.com/graphql-go/graphql"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQueries",
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
				var startPos, offset int
				startPos, _ = p.Args["start_pos"].(int)
				offset, _ = p.Args["offset"].(int)

				result, err := service.Instance().AllUsers(startPos, offset)

				if err != nil {
					return nil, err
				} else {
					return gqlUserListFromModel(result), nil
				}
			},
		},
	},
	Description: "Root Query for Picnic GraphQL Web Server",
})
