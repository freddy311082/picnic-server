package utils

import (
	"github.com/google/logger"
	"io/ioutil"
	"log"
)

func LoggerObj() *logger.Logger {
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	return logger.Init("PICNIC", true, false, ioutil.Discard)

}
