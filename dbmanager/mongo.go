package dbmanager

import (
	"context"
	"errors"
	"fmt"
	"github.com/freddy311082/picnic-server/settings"
	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	initiated     bool
}

func (dbManager *mongodbManagerImp) init() {
	// users collection
	if !dbManager.initiated {
		usersCollection := dbManager.db.Collection(utils.USERS_COLLECTION)
		usersCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		})

		dbManager.initiated = true
	}
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
		logger.Info(err.Error())
	}

	if findUser == nil { // user not found
		collection := dbManager.db.Collection(utils.USERS_COLLECTION)

		userDB := mdbUserModel{}
		userDB.initFromModel(user)
		userDB.generateNewID()
		logger.Info(userDB.ID.String())

		if result, err := collection.InsertOne(context.TODO(), userDB); err != nil {
			logger.Error(err.Error())
			return nil, err
		} else {
			user.Id = &mdbId{id: result.InsertedID.(primitive.ObjectID)}
		}

		return user, nil
	}

	msg := fmt.Sprintf("User %s already exists.", user.Email)
	logger.Error(msg)
	return nil, errors.New(msg)
}

func (dbManager *mongodbManagerImp) GetUser(email string) (*model.User, error) {
	collection := dbManager.db.Collection(utils.USERS_COLLECTION)
	query := &bson.M{
		"email": email,
	}
	user := &model.User{}
	result := collection.FindOne(context.TODO(), query)

	if result.Err() != nil {
		logger.Error(result.Err().Error())
		return nil, result.Err()
	}

	if err := result.Decode(user); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (dbManager *mongodbManagerImp) Open() error {
	var err error
	dbManager.client, err = mongo.Connect(context.TODO(), dbManager.clientOptions)

	if err == nil {
		dbManager.isOpen = true
		dbManager.db = dbManager.client.Database("picnic")
		dbManager.init()
	}

	return err
}

func (dbManager *mongodbManagerImp) AllUsers(startPosition, offset int) (model.UserList, error) {
	if startPosition < 0 {
		const msg = "start position cannot be zero or a negative number"
		logger.Error(msg)
		return nil, errors.New(msg)
	}

	findOptions := options.Find()

	if offset > 0 {
		findOptions.SetLimit(int64(offset))
	}

	if startPosition > 0 {
		findOptions.SetSkip(int64(startPosition))
	}

	collection := dbManager.db.Collection(utils.USERS_COLLECTION)
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)

	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return nil, err
	}

	var users model.UserList
	userDb := &mdbUserModel{}

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&userDb); err != nil {
			logger.Error(err.Error())
			return nil, err
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
		initiated:     false,
	}

	return manager
}

func createMongoDbManager() DBManager {
	manager := &mongodbManagerImp{
		isOpen:        false,
		client:        nil,
		clientOptions: options.Client().ApplyURI(settings.SettingsObj().DBSettingsValues().ConnectionString()),
		initiated:     false,
	}

	return manager
}
