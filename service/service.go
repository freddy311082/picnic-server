package service

import (
	"github.com/freddy311082/picnic-server/dbmanager"
	"github.com/freddy311082/picnic-server/model"
)

type Service struct {
}

func (service *Service) users(startPosition, offset int) {
}

func (service *Service) RegisterUser(user *model.User) (*model.User, error) {
	return dbmanager.DBManagerInstance().RegisterNewUser(user)
}
