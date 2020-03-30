package api

import (
	"encoding/json"
	"fmt"
	"github.com/freddy311082/picnic-server/service"
	"github.com/freddy311082/picnic-server/settings"
	"github.com/freddy311082/picnic-server/utils"
	"github.com/friendsofgo/graphiql"
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
	// Init services
	utils.PicnicLog_INFO("Starting Picnic Web Server")
	utils.PicnicLog_INFO("Initiating services...")
	if err := service.Instance().Init(); err != nil {
		utils.PicnicLog_ERROR("Error starting start the server...")
		utils.PicnicLog_ERROR(err.Error())
		utils.PicnicLog_INFO("Server stopped...")
	}
	utils.PicnicLog_INFO("Services initiated :)")

	portStr := fmt.Sprintf(":%d", settings.SettingsObj().APISettings().HttpPort())
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		utils.PicnicLog_ERROR(err.Error())
		return
	}
	http.Handle("/graphiql", graphiqlHandler)
	http.Handle("/graphql", server.getGqlHandler())

	utils.PicnicLog_INFO(settings.SettingsObj().ToString())
	http.ListenAndServe(portStr, nil)
}

func (server *gqlServerImp) getGqlHandler() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Body == nil {
			http.Error(response, "No query data.", 400)
			return
		}

		rBody := server.decodeRequest(request, response)
		if result, err := server.processQuery(rBody.Query); err != nil {
			http.Error(response, err.Error(), 400)
		} else {
			fmt.Fprintf(response, "%s", result)
		}

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

func WebServerInstance() WebServer {
	return &gqlServerImp{}
}
