package api

import (
	"encoding/json"
	"net/http"
)

type WebServer interface {
	Start()
	Stop()
}

type gqlServerImp struct {
}

func (server *gqlServerImp) processQuery(query string) string {

}

func (server *gqlServerImp) Start() {
	http.Handle("/graphql", http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Body == nil {
			return
		}

		type reqBody struct {
			Query string `json:"query"`
		}

		var rBody reqBody
		if err := json.NewDecoder(request.Body).Decode(&rBody); err != nil {
			const msg = "Error parsing JSON request body"
			http.Error(response, msg, 400)
		}

		server.processQuery(response, rBody.Query)
	}))
}

func (server *gqlServerImp) Stop() {

}
