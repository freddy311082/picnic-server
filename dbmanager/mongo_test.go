package dbmanager

import (
	"github.com/freddy311082/picnic-server/settings"
	"testing"
)

func initMongodbManagerForTesting() (*mongodbManagerImp, error) {
	settingsObj := settings.SettingsObj()
	settingsObj.DBSettingsValues().ChangeDatabase("picnic-testing")
	connString := settingsObj.DBSettingsValues().ConnectionString()
	mongodbManager := createMongoDbManagerForTesting(connString)

	if err := mongodbManager.Open(); err != nil {
		return nil, err
	}
	return mongodbManager, nil
}

func TestOpenConnectionWithCloudMongo(t *testing.T) {
	if mongodbManager, err := initMongodbManagerForTesting(); err != nil {
		t.Error(err)
	} else {
		mongodbManager.Close()
	}
}
