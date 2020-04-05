package api

import (
	"github.com/freddy311082/picnic-server/utils"
	"github.com/graphql-go/graphql"
)

func GetSchema() (*graphql.Schema, error) {
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
