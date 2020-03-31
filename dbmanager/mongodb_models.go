package dbmanager

import (
	"errors"
	"github.com/freddy311082/picnic-server/model"
	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mdbId struct {
	id primitive.ObjectID
}

func (mdbIdObj *mdbId) ToString() string {
	return mdbIdObj.id.Hex()
}

type mdbUserModel struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	LastName string             `bson:"last_name"`
	Email    string             `bson:"email"`
}

func (dbUser *mdbUserModel) generateNewID() {
	dbUser.ID = primitive.NewObjectID()
}

func (dbUser *mdbUserModel) initFromModel(user *model.User) error {
	if user == nil {
		const msg = "invalid userModel. Cannot initiate MongoDB user model from NULL"
		logger.Error(msg)
		return errors.New(msg)
	}

	if user.Id != nil {
		if objId, err := primitive.ObjectIDFromHex(user.Id.ToString()); err != nil {
			dbUser.ID = objId
		}
	}
	dbUser.Name = user.Name
	dbUser.LastName = user.LastName
	dbUser.Email = user.Email

	return nil
}

func (dbUser *mdbUserModel) toModel() *model.User {
	return &model.User{
		Id: &mdbId{
			id: dbUser.ID,
		},
		Name:     dbUser.Name,
		LastName: dbUser.LastName,
		Email:    dbUser.Email,
		Token:    "",
	}
}
