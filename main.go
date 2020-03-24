package main

import (
	"github.com/google/logger"
	"io/ioutil"
)

func main() {
	defer logger.Init("PICNIC", true, false, ioutil.Discard).Close()

	logger.Info("Starting Picnic Web Server")
	logger.Info("Started Picnic Web Server")
	logger.Info("Stopped Picnic Web Server")
}
