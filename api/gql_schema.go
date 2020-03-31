package api

import (
	"github.com/google/logger"
	"github.com/graphql-go/graphql"
)

func GetSchema() (*graphql.Schema, error) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &schema, err
}
