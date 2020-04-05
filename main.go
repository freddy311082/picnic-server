package main

import (
	"github.com/freddy311082/picnic-server/api"
	"github.com/freddy311082/picnic-server/utils"
)

func startServer() {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()
	server := api.WebServerInstance()
	server.Start()
	loggerObj.Info("Stopped Picnic Web Server")
}

func main() {
	startServer()
	//api.StartTest()
}
