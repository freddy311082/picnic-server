package service

import (
	"github.com/freddy311082/picnic-server/dbmanager"
	"github.com/freddy311082/picnic-server/model"
)

type Service interface {
	Users(startPosition, offset int) ([]model.User, error)
	RegisterUser(user *model.User) (*model.User, error)
}

type serviceImp struct {
}

func (service *serviceImp) RegisterUser(user *model.User) (*model.User, error) {
	return dbmanager.DBManagerInstance().RegisterNewUser(user)
}

func (service *serviceImp) Users(startPosition, offset int) ([]model.User, error) {
	return dbmanager.DBManagerInstance().AllUsers(startPosition, offset)
}

func ServiceInstance() Service {
	return &serviceImp{}
}
