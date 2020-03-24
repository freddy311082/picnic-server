package service

import (
	"testing"
)

func TestLoadSettingsFile(t *testing.T) {
	settingsObj := settingsImp{}

	expectedPath := "/Users/freddymartinezgarcia/go/src/github.com/freddy311082/picnic-server/service/../config/settings.json"
	resolvedPath := settingsObj.filename()

	if resolvedPath != expectedPath {
		t.Error("Filename resolved as settings.json is not correct.")
	}
}

func TestSettingsImp_LoadObject(t *testing.T) {
	settingsObj := SettingsObj()
	connStrExpected := "mongodb+srv://admin:AtN9WUWKwftoqijRPnAwT8Tam6WCJ2WjdAHwChQ8Tp0=@cluster0-uekoh.mongodb.net/picnic?retryWrites=true&w=majority"
	connStrValue := settingsObj.DBSettingsValues().ConnectionString()
	if connStrValue != connStrExpected {
		t.Error("Invalid settings data loaded.")
		t.Log("Value received: ", connStrValue)
	}
}
