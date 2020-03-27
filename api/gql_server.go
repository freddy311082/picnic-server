package api

import (
	"encoding/json"
	"fmt"
	"github.com/freddy311082/picnic-server/service"
	"github.com/freddy311082/picnic-server/utils"
	"github.com/graphql-go/graphql"
	"net/http"
)

type reqBody struct {
	Query string `json:"query"`
}

type WebServer interface {
	Start()
	Stop()
}

type gqlServerImp struct {
}

func (server *gqlServerImp) Start() {
	portStr := fmt.Sprintf(":%d", service.SettingsObj().APISettings().HttpPort())

	http.Handle("/graphql", server.getGqlHandler())
	http.ListenAndServe(portStr, nil)
}

func (server *gqlServerImp) getGqlHandler() http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Body == nil {
			return
		}

		rBody := server.decodeRequest(request, response)
		server.processQuery(rBody.Query)
	})
}

func (server *gqlServerImp) decodeRequest(request *http.Request, response http.ResponseWriter) reqBody {
	var rBody reqBody
	if err := json.NewDecoder(request.Body).Decode(&rBody); err != nil {
		const msg = "Error parsing JSON request body"
		utils.PicnicLog_ERROR(msg)
		http.Error(response, msg, 400)
	}
	return rBody
}

func (server *gqlServerImp) processQuery(query string) (string, error) {
	schema, err := GetSchema()

	if err != nil {
		return "", err
	}

	gqlParams := graphql.Params{
		Schema:        *schema,
		RequestString: query,
	}

	result := graphql.Do(gqlParams)
	responseJSON, _ := json.Marshal(result)
	return fmt.Sprintf("%s", responseJSON), nil
}

func (server *gqlServerImp) Stop() {

}
