package service

import (
	"github.com/freddy311082/picnic-server/dbmanager"
	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/utils"
)

type Service interface {
	Init() error
	Users(startPosition, offset int) ([]model.User, error)
	RegisterUser(user *model.User) (*model.User, error)
}

type serviceImp struct {
	dbManager dbmanager.DBManager
}

func (service *serviceImp) Init() error {
	loggerObj := utils.LoggerObj()
	if err := serviceInstance.dbManager.Open(); err != nil {
		loggerObj.Error(err.Error())
		panic(err.Error())
		return err
	}

	return nil
}

func (service *serviceImp) RegisterUser(user *model.User) (*model.User, error) {
	return dbmanager.DBManagerInstance().RegisterNewUser(user)
}

func (service *serviceImp) Users(startPosition, offset int) ([]model.User, error) {
	return dbmanager.DBManagerInstance().AllUsers(startPosition, offset)
}

var serviceInstance *serviceImp

func Instance() Service {
	ch := make(chan Service)

	go func() {
		if serviceInstance == nil {
			serviceInstance = &serviceImp{
				dbManager: dbmanager.DBManagerInstance(),
			}
		}

		ch <- serviceInstance
	}()

	return <-ch
}
