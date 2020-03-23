package service

import (
	"testing"
)

func TestLoadDbSettings(t *testing.T) {
	const fileContent = `{
  "db": {
    "mongodb": {
      "host": "cluster0-uekoh.mongodb.net/test?retryWrites=true&w=majority",
      "port": 27000,
      "dbname": "picnic",
      "user": "admin",
      "password": "Picnic2020"
    }
  }
}`
	settings := NewInstance()

}
