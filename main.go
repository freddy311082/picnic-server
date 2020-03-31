package main

import (
	"github.com/freddy311082/picnic-server/api"
	"github.com/freddy311082/picnic-server/utils"
	"github.com/google/logger"
)

func startServer() {
	defer utils.InitLogger()
	server := api.WebServerInstance()
	server.Start()
	logger.Info("Stopped Picnic Web Server")
}

func main() {
	startServer()
	//api.StartTest()
}
