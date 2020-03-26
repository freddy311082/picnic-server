package dbmanager

import (
	"context"
	"github.com/freddy311082/picnic-server/model"
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
	user, err := dbManager.GetUser(user.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		collection := dbManager.db.Collection("users")
		if result, err := collection.InsertOne(context.TODO(), user); err != nil {
			return nil, err
		} else {
			user.ID = result.InsertedID.(string)
		}
	}

	return user, nil
}

func (dbManager *mongodbManagerImp) GetUser(email string) (*model.User, error) {
	collection := dbManager.db.Collection("users")
	query := &bson.M{
		"email": email,
	}
	user := &model.User{}
	err := collection.FindOne(context.TODO(), query).Decode(user)

	return user, err
}

func (dbManager *mongodbManagerImp) Open() error {
	var err error
	dbManager.client, err = mongo.Connect(context.TODO(), dbManager.clientOptions)

	if err == nil {
		dbManager.isOpen = true
	}

	return err
}

func createMongoDbManager(connectionString string) *mongodbManagerImp {
	manager := &mongodbManagerImp{
		isOpen:        false,
		client:        nil,
		clientOptions: options.Client().ApplyURI(connectionString),
	}

	return manager
}
