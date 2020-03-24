package utils

import (
	"github.com/google/logger"
	"io/ioutil"
	"log"
)

func initLogger() {
	logger.Init("PICNIC", true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
}

func PicnicLog_ERROR(msg string) {
	initLogger()
	logger.Error(msg)
}

func PicnicLog_INFO(msg string) {
	initLogger()
	logger.Info(msg)
}
