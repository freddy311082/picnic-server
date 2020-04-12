package api

import (
	"errors"
	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/service"
	"github.com/freddy311082/picnic-server/utils"
	"github.com/graphql-go/graphql"
	"time"
)

type gqlUserRsp struct {
	ID       string
	Name     string
	LastName string
	Email    string
}

type gqlUserListResponse []*gqlUserRsp

func gqlUserFromModel(user *model.User) *gqlUserRsp {
	return &gqlUserRsp{
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

type gqlProjectRsp struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time `json:"created_at"`
	Owner       *gqlUserRsp
	Customer    *gqlCustomerRsp
}

func (obj *gqlProjectRsp) initFromModel(project *model.Project) {
	obj.ID = project.ID.ToString()
	obj.Name = project.Name
	obj.Description = project.Description
	obj.CreatedAt = project.CreatedAt
}

type gqlProjectListRsp []gqlProjectRsp

type gqlCustomerRsp struct {
	ID       string
	Name     string
	Cuit     string
	Projects *gqlProjectRsp
}

func GetSchema() (*graphql.Schema, error) {
	UserType := graphql.NewObject(graphql.ObjectConfig{
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

	ProjectType := graphql.NewObject(graphql.ObjectConfig{
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
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					return nil, nil
				},
			},
		},
		Description: "Project object definition",
	})

	CustomerType := graphql.NewObject(graphql.ObjectConfig{
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
				Type: &graphql.List{
					OfType: ProjectType,
				},
				Description: "List of project linked to this Customer.",
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					if projects, ok := p.Args["projects"].([]interface{}); ok {
						var ids model.IDList
						for _, projectId := range projects {
							if id, ok := projectId.(string); ok {
								ids = append(ids, service.Instance().NewIDFromString(id))
							} else {
								return nil, errors.New("invalid project id")
							}
						}

						return service.Instance().AllProjectWhereIDIsIn(ids)
					}

					return model.IDList{}, nil
				},
			},
		},
	})

	CustomerType.AddFieldConfig("customer", &graphql.Field{
		Type: CustomerType,
	})

	// Queries
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

	// Mutations
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

	// Schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err != nil {
		utils.LoggerObj().Error(err.Error())
		return nil, err
	}

	return &schema, err
}
