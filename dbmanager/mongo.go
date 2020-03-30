package dbmanager

import (
	"context"
	"fmt"
	"github.com/freddy311082/picnic-server/settings"

	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbManagerImp struct {
	isOpen        bool
	client        *mongo.Client
	clientOptions *options.ClientOptions
	db            *mongo.Database
}

func (dbManager *mongodbManagerImp) init() {
	// users collection
	usersCollection := dbManager.db.Collection(utils.USERS_COLLECTION)
	usersCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
}

func (dbManager *mongodbManagerImp) Close() error {
	var err error
	if dbManager.isOpen {
		err = dbManager.client.Disconnect(context.TODO())
		dbManager.isOpen = false
		dbManager.client = nil
		dbManager.db = nil
	}

	return err
}

func (dbManager *mongodbManagerImp) IsOpen() bool {
	return dbManager.isOpen
}

func (dbManager *mongodbManagerImp) RegisterNewUser(user *model.User) (*model.User, error) {
	findUser, err := dbManager.GetUser(user.Email)

	if err != nil {
		utils.PicnicLog_INFO(err.Error())
	}

	if findUser == nil { // user not found
		collection := dbManager.db.Collection(utils.USERS_COLLECTION)
		if result, err := collection.InsertOne(context.TODO(), user); err != nil {
			return nil, utils.ErrorAndLog(err.Error())
		} else {
			user.ID = result.InsertedID.(string)
		}

		return user, nil
	}

	return nil, utils.ErrorAndLog(fmt.Sprintf("User %s already exists.", user.Email))
}

func (dbManager *mongodbManagerImp) GetUser(email string) (*model.User, error) {
	collection := dbManager.db.Collection(utils.USERS_COLLECTION)
	query := &bson.M{
		"email": email,
	}
	user := &model.User{}
	result := collection.FindOne(context.TODO(), query)

	if result.Err() != nil {
		utils.PicnicLog_ERROR(result.Err().Error())
		return nil, result.Err()
	}

	if err := result.Decode(user); err != nil {
		return nil, utils.ErrorAndLog(err.Error())
	}

	return user, nil
}

func (dbManager *mongodbManagerImp) Open() error {
	var err error
	dbManager.client, err = mongo.Connect(context.TODO(), dbManager.clientOptions)

	if err == nil {
		dbManager.isOpen = true
		dbManager.db = dbManager.client.Database("picnic")
	}

	return err
}

func (dbManager *mongodbManagerImp) AllUsers(startPosition, offset int) (model.UserList, error) {
	if startPosition < 0 {
		return nil, utils.ErrorAndLog("start position cannot be zero or a negative number")
	}

	findOptions := options.Find()

	if offset > 0 {
		findOptions.SetLimit(int64(offset))
	}

	if startPosition > 0 {
		findOptions.SetSkip(int64(startPosition))
	}

	collection := dbManager.db.Collection(utils.USERS_COLLECTION)
	cursor, err := collection.Find(context.TODO(), nil, findOptions)

	if err != nil {
		utils.PicnicLog_ERROR(fmt.Sprintf("%s", err))
		return nil, err
	}

	var users model.UserList
	var userDb *mdbUserModel

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(userDb); err != nil {
			return nil, utils.ErrorAndLog(err.Error())
		}

		user := userDb.toModel()
		users = append(users, *user)
	}

	return users, nil
}

func createMongoDbManagerForTesting(connectionString string) *mongodbManagerImp {
	manager := &mongodbManagerImp{
		isOpen:        false,
		client:        nil,
		clientOptions: options.Client().ApplyURI(connectionString),
	}

	return manager
}

func createMongoDbManager() DBManager {
	manager := &mongodbManagerImp{
		isOpen:        false,
		client:        nil,
		clientOptions: options.Client().ApplyURI(settings.SettingsObj().DBSettingsValues().ConnectionString()),
	}

	return manager
}
