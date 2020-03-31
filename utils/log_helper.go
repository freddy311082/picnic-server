package utils

import (
	"github.com/google/logger"
	"io/ioutil"
	"log"
)

func InitLogger() {
	logger.Init("PICNIC", true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
}
