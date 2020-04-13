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

type gqlUserListRsp []*gqlUserRsp

func gqlUserFromModel(user *model.User) *gqlUserRsp {
	return &gqlUserRsp{
		ID:       user.ID.ToString(),
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
	}
}

func gqlUserListFromModel(userList model.UserList) gqlUserListRsp {
	var gqlUserList gqlUserListRsp

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

type gqlProjectListRsp []*gqlProjectRsp

func gqlProjectFromModel(project *model.Project) *gqlProjectRsp {
	result := &gqlProjectRsp{
		ID:          project.ID.ToString(),
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
	}

	if project.Owner != nil {
		result.Owner = gqlUserFromModel(project.Owner)
	}

	if project.Customer != nil {
		result.Customer = gqlCustomerFromModel(project.Customer)
	}

	return result
}

func gqlProjectListFromModel(projects model.ProjectList) gqlProjectListRsp {
	var result gqlProjectListRsp

	for _, project := range projects {
		result = append(result, gqlProjectFromModel(project))
	}

	return result
}

type gqlCustomerRsp struct {
	ID       string
	Name     string
	Cuit     string
	Projects gqlProjectListRsp
}

type gqlCustomerListRsp []*gqlCustomerRsp

func gqlCustomerFromModel(customer *model.Customer) *gqlCustomerRsp {
	return &gqlCustomerRsp{
		ID:       customer.ID.ToString(),
		Name:     customer.Name,
		Cuit:     customer.Cuit,
		Projects: gqlProjectListFromModel(customer.Projects),
	}
}

func gqlCustomerListFromModel(customers model.CustomerList) gqlCustomerListRsp {
	var result gqlCustomerListRsp

	for _, customer := range customers {
		result = append(result, gqlCustomerFromModel(customer))
	}

	return result
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

	ProjectType.AddFieldConfig("customer", &graphql.Field{
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
			"customers": &graphql.Field{
				Type: &graphql.List{
					OfType: CustomerType,
				},
				Args: nil,
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					if customers, err := service.Instance().AllCustomers(); err != nil {
						return nil, err
					} else {
						return gqlCustomerListFromModel(customers), nil
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
					"owner": &graphql.ArgumentConfig{
						Type:        &graphql.NonNull{OfType: graphql.ID},
						Description: "User ID corresponding to the project owner which is the one who created this project.",
					},
				},
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					var name, description string
					var id model.ID

					name, _ = p.Args["name"].(string)
					if name == "" {
						return nil, errors.New("name is required and cannot be empty")
					}

					if value, ok := p.Args["description"].(string); ok {
						description = value
					} else {
						description = ""
					}

					if value, ok := p.Args["owner"].(string); ok {
						id = service.Instance().NewIDFromString(value)
					} else {
						return nil, errors.New("owner id cannot be nil")
					}

					if result, err := service.Instance().CreateProject(&model.Project{
						Name:        name,
						Description: description,
						CreatedAt:   time.Now(),
						Owner:       &model.User{ID: id},
					}); err != nil {
						return nil, err
					} else {
						return gqlProjectFromModel(result), nil
					}
				},
			},
			"createCustomer": &graphql.Field{
				Type: CustomerType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: &graphql.NonNull{OfType: graphql.String},
					},
					"cuit": &graphql.ArgumentConfig{
						Type: &graphql.NonNull{OfType: graphql.String},
					},
				},
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					var name, cuit string

					if value, ok := p.Args["name"].(string); !ok {
						return nil, errors.New("invalid name value")
					} else {
						name = value
					}

					if value, ok := p.Args["cuit"].(string); !ok {
						return nil, errors.New("")
					} else {
						cuit = value
					}

					if result, err := service.Instance().CreateCustomer(&model.Customer{
						Name: name,
						Cuit: cuit,
					}); err != nil {
						return nil, err
					} else {
						return gqlCustomerFromModel(result), nil
					}
				},
				Description: "Create new Customer in the system",
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
