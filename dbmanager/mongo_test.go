package dbmanager

import (
	"github.com/freddy311082/picnic-server/service"
	"testing"
)

func initMongodbManagerForTesting() (*mongodbManagerImp, error) {
	settingsObj := service.SettingsObj()
	settingsObj.DBSettingsValues().ChangeDatabase("picnic-testing")
	connString := settingsObj.DBSettingsValues().ConnectionString()
	mongodbManager := createMongoDbManager(connString)

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
