package dbmanager

import (
	"errors"
	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/utils"
)

type mdbUserModel struct {
	id       string `json:"_id"`
	name     string `json:"name"`
	lastName string `json:"last_name"`
	email    string `json:"email"`
}

func (dbUser *mdbUserModel) initFromModel(user *model.User) error {
	if user == nil {
		const msg = "invalid userModel. Cannot initiate MongoDB user model from NULL"
		utils.PicnicLog_ERROR(msg)
		return errors.New(msg)
	}

	dbUser.id = user.ID
	dbUser.name = user.Name
	dbUser.lastName = user.LastName
	dbUser.email = user.Email

	return nil
}

func (dbUser *mdbUserModel) toModel() *model.User {
	return &model.User{
		ID:       dbUser.id,
		Name:     dbUser.name,
		LastName: dbUser.lastName,
		Email:    dbUser.email,
		Token:    "",
	}
}
