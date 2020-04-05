package api

import (
	"encoding/json"
	"fmt"
	"github.com/freddy311082/picnic-server/service"
	"github.com/freddy311082/picnic-server/settings"
	"github.com/freddy311082/picnic-server/utils"
	"github.com/friendsofgo/graphiql"
	"github.com/go-chi/chi"
	"github.com/graphql-go/graphql"
	"github.com/rs/cors"
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
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()
	// Init services
	loggerObj.Info("Starting Picnic Web Server")
	loggerObj.Info("Initiating services...")
	if err := service.Instance().Init(); err != nil {
		loggerObj.Error("Error starting start the server...")
		loggerObj.Error(err.Error())
		loggerObj.Info("Server stopped...")
	}
	loggerObj.Info("Services initiated :)")

	portStr := fmt.Sprintf(":%d", settings.SettingsObj().APISettings().HttpPort())
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		loggerObj.Error(err.Error())
		return
	}
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	router.Handle("/graphiql", graphiqlHandler)
	router.Handle("/graphql", server.getGqlHandler())

	loggerObj.Info(settings.SettingsObj().ToString())
	http.ListenAndServe(portStr, router)
}

func (server *gqlServerImp) getGqlHandler() http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		loggerObj := utils.LoggerObj()
		if request.Body == nil {
			const msg = "Error 400: No query data."
			loggerObj.Error(msg)
			http.Error(response, msg, 400)
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
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if err := json.NewDecoder(request.Body).Decode(&rBody); err != nil {
		const msg = "Error parsing JSON request body"
		loggerObj.Error(msg)
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
