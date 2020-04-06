package dbmanager

import (
	"errors"
	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/utils"
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

func (dbUser *mdbUserModel) initFromModel(user *model.User) error {
	if user == nil {
		const msg = "invalid userModel. Cannot initiate MongoDB user model from NULL object"
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		loggerObj.Error(msg)

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

type mdbProjectModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
	OwnerID     primitive.ObjectID `bson:"owner_id"`
}

func (dbProject *mdbProjectModel) initFromModel(project *model.Project) error {
	if project == nil {
		const msg = "invalid projectModel. Cannot initialize MongoDB project model from NULL object"
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		loggerObj.Error(msg)
		return errors.New(msg)
	}

	if project.ID != nil {
		if objId, err := primitive.ObjectIDFromHex(project.ID.ToString()); err != nil {
			dbProject.ID = objId
		}
	}

	dbProject.Name = project.Name
	dbProject.Description = project.Description
	dbProject.CreatedAt = primitive.NewDateTimeFromTime(project.CreatedAt)

	return nil
}

func (dbProject *mdbProjectModel) toModel() *model.Project {
	return &model.Project{
		ID:          &mdbId{id: dbProject.ID},
		Name:        dbProject.Name,
		Description: dbProject.Description,
		CreatedAt:   dbProject.CreatedAt.Time(),
		Owner:       nil,
		Customer:    nil,
		Fields:      nil,
	}
}
