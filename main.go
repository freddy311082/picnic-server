package main

import (
	"github.com/freddy311082/picnic-server/api"
	"github.com/google/logger"
	"io/ioutil"
)

func startServer() {
	defer logger.Init("PICNIC", true, false, ioutil.Discard).Close()
	server := api.WebServerInstance()
	server.Start()
	logger.Info("Stopped Picnic Web Server")
}

func main() {
	startServer()
	//api.StartTest()
}
