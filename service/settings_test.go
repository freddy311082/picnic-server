package service

import (
	"testing"
)

type mockSettings struct {
	settingsImp
}

func (settings *mockSettings) fileContent() ([]byte, error) {
	return []byte(`{
  "db": {
    "mongodb": {
      "host": "cluster0-uekoh.mongodb.net/test?retryWrites=true&w=majority",
      "port": 27000,
      "dbname": "picnic",
      "user": "admin",
      "password": "Picnic2020"
    }
  }
}`), nil
}

func TestLoadSettingsFile(t *testing.T) {
	settingsObj := settingsImp{}

	expectedPath := "/Users/freddymartinezgarcia/go/src/github.com/freddy311082/picnic-server/service/../config/settings.json"
	resolvedPath := settingsObj.filename()

	if resolvedPath != expectedPath {
		t.Error("Filename resolved as settings.json is not correct.")
	}
}
