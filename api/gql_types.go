package api

import (
	"errors"
	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/service"
	"github.com/graphql-go/graphql"
	"time"
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
		ID:       user.ID.ToString(),
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
	}
}

func gqlUserListFromModel(userList model.UserList) gqlUserListResponse {
	var gqlUserList gqlUserListResponse

	for _, user := range userList {
		gqlUserList = append(gqlUserList, gqlUserFromModel(user))
	}

	return gqlUserList
}

var ProjectType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProjectType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"owner": &graphql.Field{
			Type: UserType,
		},
		"customer": &graphql.Field{
			Type: Customer,
		},
	},
	Description: "Project object definition",
})

type gqlProjectTypeRsp struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time `json:"created_at"`
	Owner       *gqlUserResponse
	Customer    *gqlCustomer
}

var Customer = graphql.NewObject(graphql.ObjectConfig{
	Name: "CustomerType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"cuit": &graphql.Field{
			Type: graphql.String,
		},
		"projects": &graphql.Field{
			Type:        &graphql.List{OfType: ProjectType},
			Description: "List of project linked to this Customer.",
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				if projects, ok := p.Args["projects"].([]interface{}); ok {

					var ids model.IDList
					for _, projectId := range projects {
						if id, ok := projectId.(string); ok {
							ids = append(ids, service.Instance().NewIDFromString(id))
						} else {
							return nil, errors.New("error: project is not valid")
						}
					}

					return ids, nil
				}

				return model.IDList{}, errors.New("invalid project id list")
			},
		},
	},
})

type gqlCustomer struct {
	ID   string
	Name string
}
