package dbmanager

import (
	"context"
	"errors"
	"fmt"
	"github.com/freddy311082/picnic-server/settings"
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

func (dbManager *mongodbManagerImp) collection(name string) *mongo.Collection {
	return dbManager.db.Collection(name)
}

func (dbManager *mongodbManagerImp) AllProjects(startPosition, offset int) (model.ProjectList, error) {
	panic("implement me")
}

func (dbManager *mongodbManagerImp) CreateProject(project *model.Project) (*model.Project, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	projectDb := &mdbProjectModel{}
	if err := projectDb.initFromModel(project); err != nil {
		return nil, err
	}

	collection := dbManager.collection(utils.PROJECTS_COLLECTION)
	if result, err := collection.InsertOne(context.TODO(), projectDb); err != nil {
		loggerObj.Error(err.Error())
		return nil, err
	} else {
		project.ID = &mdbId{id: result.InsertedID.(primitive.ObjectID)}
	}

	return project, nil
}

func (dbManager *mongodbManagerImp) GetProject(projectId model.ID) (*model.Project, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	dbId, err := primitive.ObjectIDFromHex(projectId.ToString())
	if err != nil {
		loggerObj.Error(err)
		return nil, err
	}

	collection := dbManager.collection(utils.PROJECTS_COLLECTION)

	result := collection.FindOne(context.TODO(), bson.M{utils.PROJECT_ID_FIELD: dbId})
	return dbManager.decodeBsonIntoProjectModel(result)
}

func (dbManager *mongodbManagerImp) decodeBsonIntoProjectModel(result *mongo.SingleResult) (*model.Project, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if result.Err() != nil {
		loggerObj.Error(result.Err())
		return nil, result.Err()
	}

	projectDb := &mdbProjectModel{}
	if err := result.Decode(projectDb); err != nil {
		loggerObj.Error(err)
		return nil, err
	}

	project := projectDb.toModel()
	if owner, err := dbManager.GetUserByID(&mdbId{id: projectDb.OwnerID}); err != nil {
		return nil, err
	} else {
		project.Owner = owner
	}

	return project, nil
}

func (dbManager *mongodbManagerImp) UpdateProject(project *model.Project) (*model.Project, error) {
	panic("implement me")
}

func (dbManager *mongodbManagerImp) init() {
	// users collection
	if !dbManager.initiated {
		usersCollection := dbManager.collection(utils.USERS_COLLECTION)
		usersCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
			Keys:    bson.M{utils.USER_EMAIL_FIELD: 1},
			Options: options.Index().SetUnique(true),
		})
		usersCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
			Keys:    bson.M{utils.PROJECT_NAME_FIELD: 1, utils.PROJECT_OWNER_ID_FIELD: 1},
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
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()
	findUser, err := dbManager.GetUserByEmail(user.Email)

	if err != nil {
		loggerObj.Info(err)
	}

	if findUser == nil { // user not found
		collection := dbManager.collection(utils.USERS_COLLECTION)

		userDB := mdbUserModel{}
		userDB.initFromModel(user)
		userDB.ID = primitive.NewObjectID()
		loggerObj.Info(userDB.ID.String())

		if result, err := collection.InsertOne(context.TODO(), userDB); err != nil {
			loggerObj.Error(err.Error())
			return nil, err
		} else {
			user.Id = &mdbId{id: result.InsertedID.(primitive.ObjectID)}
		}

		return user, nil
	}

	msg := fmt.Sprintf("User %s already exists.", user.Email)
	loggerObj.Error(msg)
	return nil, errors.New(msg)
}

func (dbManager *mongodbManagerImp) GetUserByEmail(email string) (*model.User, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	collection := dbManager.collection(utils.USERS_COLLECTION)
	query := &bson.M{
		utils.USER_EMAIL_FIELD: email,
	}

	result := collection.FindOne(context.TODO(), query)

	if result.Err() != nil {
		loggerObj.Error(result.Err().Error())
		return nil, result.Err()
	}

	return dbManager.decodeBsonIntoUserModel(result)
}

func (dbManager *mongodbManagerImp) GetUserByID(id model.ID) (*model.User, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	dbId, err := primitive.ObjectIDFromHex(id.ToString())

	if err != nil {
		loggerObj.Error(err)
		return nil, err
	}

	collection := dbManager.collection(utils.USERS_COLLECTION)
	result := collection.FindOne(context.TODO(), bson.M{utils.USER_ID_FIELD: dbId})

	if result.Err() != nil {
		loggerObj.Error(err)
		return nil, err
	}

	return dbManager.decodeBsonIntoUserModel(result)
}

func (dbManager *mongodbManagerImp) decodeBsonIntoUserModel(result *mongo.SingleResult) (*model.User, error) {
	panic("Must be implemented")
}

func (dbManager *mongodbManagerImp) Open() error {
	var err error
	dbManager.client, err = mongo.Connect(context.TODO(), dbManager.clientOptions)

	if err == nil {
		dbManager.isOpen = true
		dbManager.db = dbManager.client.Database(settings.SettingsObj().DBSettingsValues().DbName())
		dbManager.init()
	}

	return err
}

func (dbManager *mongodbManagerImp) AllUsers(startPosition, offset int) (model.UserList, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if startPosition < 0 {
		const msg = "start position cannot be zero or a negative number"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	findOptions := options.Find()

	if offset > 0 {
		findOptions.SetLimit(int64(offset))
	}

	if startPosition > 0 {
		findOptions.SetSkip(int64(startPosition))
	}

	collection := dbManager.collection(utils.USERS_COLLECTION)
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)

	if err != nil {
		loggerObj.Error(fmt.Sprintf("%s", err))
		return nil, err
	}

	var users model.UserList
	userDb := &mdbUserModel{}

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&userDb); err != nil {
			loggerObj.Error(err.Error())
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
