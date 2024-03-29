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

type cacheKey struct {
	id             primitive.ObjectID
	collectionName string
}

type mongodbManagerImp struct {
	isOpen        bool
	client        *mongo.Client
	clientOptions *options.ClientOptions
	db            *mongo.Database
	initiated     bool
	cache         map[cacheKey]interface{}
}

func (dbManager *mongodbManagerImp) AllProjectsFromCustomer(customerId model.ID) (model.ProjectList, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if dbCustomerId, err := dbManager.modelIDtoMongoID(customerId, loggerObj); err != nil {
		return model.ProjectList{}, err
	} else {
		collection := dbManager.collection(utils.PROJECTS_COLLECTION)

		if cursor, err := collection.Find(context.TODO(), bson.M{utils.PROJECT_CUSTOMER_ID_FIELD: dbCustomerId}); err != nil {
			loggerObj.Error(err)
			return model.ProjectList{}, err
		} else {
			return dbManager.decodeBsonIntoProjectListModel(cursor, loggerObj)
		}
	}
}

func (dbManager *mongodbManagerImp) GetOwnerFromProjectID(projectId model.ID) (*model.User, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()
	var err error
	var dbProjectId primitive.ObjectID
	var user *model.User

	if dbProjectId, err = dbManager.modelIDtoMongoID(projectId, loggerObj); err != nil {
		return nil, err
	} else if user, err = dbManager.GetUserByID(&mdbId{id: dbProjectId}); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (dbManager *mongodbManagerImp) GetCustomerByID(customerId model.ID) (*model.Customer, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if id, err := dbManager.modelIDtoMongoID(customerId, loggerObj); err != nil {
		return nil, err
	} else {
		collection := dbManager.collection(utils.CUSTOMERS_COLLECTION)

		if result := collection.FindOne(context.TODO(), bson.M{utils.CUSTOMER_ID_FIELD: id}); result.Err() != nil {
			loggerObj.Error(err)
			return nil, result.Err()
		} else {
			return dbManager.decodeBsonIntoCustomerModel(result, loggerObj)
		}
	}
}

func (dbManager *mongodbManagerImp) existsObject(id primitive.ObjectID, collectionName string) (bool, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	collection := dbManager.collection(collectionName)
	if count, err := collection.CountDocuments(context.TODO(), bson.M{"_id": id}); err != nil {
		loggerObj.Error(err)
		return false, err
	} else {
		return count > 0, nil
	}
}

func (dbManager *mongodbManagerImp) AllUsersWhereIDIsIn(ids model.IDList) (model.UserList, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if dbIds, err := dbManager.modelIDsToMongoIDs(ids, loggerObj); err != nil {
		loggerObj.Error(err)
		return model.UserList{}, err
	} else if cursor, queryErr := dbManager.collection(utils.USERS_COLLECTION).Find(context.TODO(),
		bson.M{
			utils.USER_ID_FIELD: bson.M{"$in": dbIds},
		}); queryErr != nil {
		loggerObj.Error(queryErr)
		return model.UserList{}, err
	} else {
		return dbManager.decodeBsonIntoUserListModel(cursor, loggerObj)
	}
}

func (dbManager *mongodbManagerImp) AllCustomersWhereIDIsIn(ids model.IDList) (model.CustomerList, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if mbdIds, err := dbManager.modelIDsToMongoIDs(ids, loggerObj); err != nil {
		return model.CustomerList{}, err
	} else {
		collection := dbManager.collection(utils.CUSTOMERS_COLLECTION)

		if cursor, err := collection.Find(context.TODO(),
			bson.M{
				"_id": bson.M{"$in": mbdIds},
			}); err != nil {
			loggerObj.Error(err)
			return model.CustomerList{}, nil
		} else {
			return dbManager.decodeBsonIntoCustomerListModel(cursor, loggerObj)
		}
	}
}

func (dbManager *mongodbManagerImp) AllProjectWhereIDIsIn(ids model.IDList) (model.ProjectList, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if mdbIds, err := dbManager.modelIDsToMongoIDs(ids, loggerObj); err != nil {
		return model.ProjectList{}, err
	} else {
		collection := dbManager.collection(utils.PROJECTS_COLLECTION)

		if cursor, err := collection.Find(context.TODO(),
			bson.M{
				"_id": bson.M{"$in": mdbIds},
			}); err != nil {
			loggerObj.Error(err)
			return model.ProjectList{}, nil
		} else {
			return dbManager.decodeBsonIntoProjectListModel(cursor, loggerObj)
		}
	}
}

func (dbManager *mongodbManagerImp) decodeBsonIntoCustomerModel(
	singleResult *mongo.SingleResult,
	loggerObj *logger.Logger) (*model.Customer, error) {

	customerDb := &mdbCustomerModel{}
	if err := singleResult.Decode(customerDb); err != nil {
		loggerObj.Error(err)
		return nil, err
	} else {
		return customerDb.toModel()
	}
}

func (dbManager *mongodbManagerImp) decodeBsonIntoCustomerListModel(
	cursor *mongo.Cursor,
	loggerObj *logger.Logger) (model.CustomerList, error) {

	var result model.CustomerList

	for cursor.Next(context.TODO()) {
		customerDb := &mdbCustomerModel{}
		if err := cursor.Decode(&customerDb); err != nil {
			loggerObj.Error(err)
			return nil, err
		} else if customer, decodeErr := customerDb.toModel(); decodeErr != nil {
			return nil, err
		} else {
			result = append(result, customer)
		}
	}

	return result, nil
}

func (dbManager *mongodbManagerImp) AllCustomers() (model.CustomerList, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()
	collection := dbManager.collection(utils.CUSTOMERS_COLLECTION)

	if cursor, err := collection.Find(context.TODO(), bson.M{}); err != nil {
		loggerObj.Error(err)
		return model.CustomerList{}, nil
	} else {
		return dbManager.decodeBsonIntoCustomerListModel(cursor, loggerObj)
	}
}

func (dbManager *mongodbManagerImp) CreateCustomer(customer *model.Customer) (*model.Customer, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	customerDb := &mdbCustomerModel{}
	customerDb.initFromModel(customer)
	customerDb.ID = primitive.NewObjectID()

	collection := dbManager.collection(utils.CUSTOMERS_COLLECTION)
	if result, err := collection.InsertOne(context.TODO(), customerDb); err != nil {
		loggerObj.Error(err)
		return nil, err
	} else {
		customer.ID = dbManager.mongoIdToModelID(result.InsertedID.(primitive.ObjectID))
		return customer, nil
	}
}

func (dbManager *mongodbManagerImp) UpdateCustomer(customer *model.Customer) (*model.Customer, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	collection := dbManager.collection(utils.CUSTOMERS_COLLECTION)

	customerDb := &mdbCustomerModel{}
	customerDb.initFromModel(customer)
	if result, err := collection.UpdateOne(context.TODO(),
		bson.D{{utils.CUSTOMER_ID_FIELD, customerDb.ID}},
		customerDb); err != nil {
		loggerObj.Error(err)
		return nil, err
	} else if result.MatchedCount != 1 {
		var msg = fmt.Sprintf("nothing to update. CustomerID (%s) was not found", customer.ID.ToString())
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	} else {
		return customer, nil
	}
}

func (dbManager *mongodbManagerImp) DeleteCustomer(customerId model.ID) error {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	collection := dbManager.collection(utils.CUSTOMERS_COLLECTION)
	if dbId, err := primitive.ObjectIDFromHex(customerId.ToString()); err != nil {
		loggerObj.Error(err)
		return err
	} else {
		var result *mongo.DeleteResult
		if result, err = collection.DeleteOne(context.TODO(), bson.M{utils.CUSTOMER_ID_FIELD: dbId}); err != nil {
			logger.Error(err)
			return err
		} else if result.DeletedCount == 0 {
			var msg = fmt.Sprintf("customer id %s not found", customerId.ToString())
			loggerObj.Error(msg)
			return errors.New(msg)
		} else {
			return nil
		}
	}
}

func (dbManager *mongodbManagerImp) DeleteCustomers(ids model.IDList) error {
	panic("implement me")
}

func (dbManager *mongodbManagerImp) DeleteProjects(ids model.IDList) error {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if mongoIds, err := dbManager.modelIDsToMongoIDs(ids, loggerObj); err != nil {
		loggerObj.Error(err)
		return err
	} else {
		collection := dbManager.collection(utils.PROJECTS_COLLECTION)
		if result, err := collection.DeleteMany(
			context.TODO(),
			bson.M{
				"$in": mongoIds,
			}); err != nil {
			loggerObj.Error(err)
			return err
		} else {
			loggerObj.Infof("Deleted %d projects.", result.DeletedCount)
		}
	}

	return nil
}

func (dbManager *mongodbManagerImp) mongoIDsToModelIDs(mdbIds []primitive.ObjectID) model.IDList {
	if len(mdbIds) > 0 {
		var result model.IDList

		for _, id := range mdbIds {
			result = append(result, &mdbId{id: id})
		}

		return result
	} else {
		return make(model.IDList, 0)
	}
}

func (dbManager *mongodbManagerImp) modelIDtoMongoID(
	id model.ID,
	loggerObj *logger.Logger) (primitive.ObjectID, error) {

	if dbId, err := primitive.ObjectIDFromHex(id.ToString()); err != nil {
		loggerObj.Error(err)
		return primitive.ObjectID{}, err
	} else {
		return dbId, nil
	}
}

func (dbManager *mongodbManagerImp) modelIDsToMongoIDs(
	ids model.IDList,
	loggerObj *logger.Logger) ([]primitive.ObjectID, error) {

	if ids == nil || len(ids) == 0 {
		return make([]primitive.ObjectID, 0), nil
	}

	var mongoIds []primitive.ObjectID

	for _, modelId := range ids {
		if id, err := dbManager.modelIDtoMongoID(modelId, loggerObj); err != nil {
			loggerObj.Error(err)
			return []primitive.ObjectID{}, err
		} else {
			mongoIds = append(mongoIds, id)
		}
	}

	return mongoIds, nil
}

func (dbManager *mongodbManagerImp) DeleteProject(projectId model.ID) error {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if id, err := primitive.ObjectIDFromHex(projectId.ToString()); err != nil {
		loggerObj.Error(err)
		return err
	} else {
		collection := dbManager.collection(utils.PROJECTS_COLLECTION)
		if _, errDelete := collection.DeleteOne(context.TODO(), bson.M{utils.PROJECT_ID_FIELD: id}); errDelete != nil {
			loggerObj.Error(err)
			return err
		}
	}

	return nil
}

func (dbManager *mongodbManagerImp) getMongoUserID(user *model.User) (*primitive.ObjectID, error) {
	var userId primitive.ObjectID
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if user.ID == nil {
		if userObj, err := dbManager.GetUserByEmail(user.Email); err != nil {
			return nil, err
		} else {
			userId, _ = primitive.ObjectIDFromHex(userObj.ID.ToString())
		}

	} else if id, err := primitive.ObjectIDFromHex(user.ID.ToString()); err != nil {
		const msg = "invalid user id. Error converting % to MongoDB id format"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	} else {
		userId = id
	}

	return &userId, nil
}

func (dbManager *mongodbManagerImp) AllProjectFromUser(user *model.User) (model.ProjectList, error) {
	var ownerId primitive.ObjectID
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if userId, err := dbManager.getMongoUserID(user); err != nil {
		return model.ProjectList{}, err
	} else {
		ownerId = *userId
	}

	collection := dbManager.collection(utils.PROJECTS_COLLECTION)

	if cursor, err := collection.Find(context.TODO(), bson.M{utils.PROJECT_OWNER_ID_FIELD: ownerId}); err != nil {
		loggerObj.Error(err)
		return model.ProjectList{}, err
	} else {
		return dbManager.decodeBsonIntoProjectListModel(cursor, loggerObj)
	}
}

func (dbManager *mongodbManagerImp) DeleteUser(email string) error {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	collection := dbManager.collection(utils.USERS_COLLECTION)

	if _, err := collection.DeleteOne(context.TODO(), bson.M{utils.USER_EMAIL_FIELD: email}); err != nil {
		loggerObj.Errorf("Error deleting user: %s. Error message: %s", email, err.Error())
		return err
	} else {
		loggerObj.Infof("Deleted user %s", email)
	}

	return nil
}

func (dbManager *mongodbManagerImp) collection(name string) *mongo.Collection {
	return dbManager.db.Collection(name)
}

func (dbManager *mongodbManagerImp) AllProjects(startPosition, offset int) (model.ProjectList, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	findOptions := options.Find()

	if offset > 0 {
		findOptions.SetLimit(int64(offset))
	}

	if startPosition > 0 {
		findOptions.SetLimit(int64(startPosition))
	}

	collection := dbManager.collection(utils.PROJECTS_COLLECTION)
	if cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions); err != nil {
		loggerObj.Error(err)
		return nil, err
	} else {
		return dbManager.decodeBsonIntoProjectListModel(cursor, loggerObj)
	}
}

func (dbManager *mongodbManagerImp) decodeBsonIntoProjectListModel(
	cursor *mongo.Cursor,
	loggerObj *logger.Logger) (model.ProjectList, error) {

	var projects model.ProjectList
	ownerCache := map[primitive.ObjectID]*model.User{}

	for cursor.Next(context.TODO()) {
		projectDb := &mdbProjectModel{}
		if err := cursor.Decode(&projectDb); err != nil {
			loggerObj.Error(err)
			return nil, err
		}

		err := dbManager.updateUserCache(ownerCache, projectDb, loggerObj)
		if err != nil {
			return model.ProjectList{}, err
		}

		project := projectDb.toModel()
		projects = append(projects, project)
	}

	return projects, nil
}

func (dbManager *mongodbManagerImp) updateUserCache(
	userCache map[primitive.ObjectID]*model.User,
	projectDb *mdbProjectModel,
	loggerObj *logger.Logger) error {

	if _, ok := userCache[projectDb.OwnerID]; !ok {
		if owner, err := dbManager.GetUserByID(&mdbId{id: projectDb.OwnerID}); err != nil {
			loggerObj.Error(err)
			return err
		} else {
			userCache[projectDb.OwnerID] = owner
		}

	}
	return nil
}

func (dbManager *mongodbManagerImp) CreateProject(project *model.Project) (*model.Project, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	projectDb := &mdbProjectModel{}
	if err := projectDb.initFromModel(project); err != nil {
		return nil, err
	}

	if exists, err := dbManager.existsObject(projectDb.CustomerID, utils.CUSTOMERS_COLLECTION); err != nil {
		return nil, err
	} else if !exists {
		err := errors.New("cannot create project with invalid customer id. CustomerID id doesn't exists")
		loggerObj.Error(err)
		return nil, err
	}

	collection := dbManager.collection(utils.PROJECTS_COLLECTION)
	projectDb.ID = primitive.NewObjectID()
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
	return project, nil
}

func (dbManager *mongodbManagerImp) UpdateProject(project *model.Project) (*model.Project, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if project == nil || project.ID == nil {
		const msg = "invalid project. Neither project object nor project ID can be NULL"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	collection := dbManager.collection(utils.PROJECTS_COLLECTION)
	projectDb := &mdbProjectModel{}

	if err := projectDb.initFromModel(project); err != nil {
		return nil, err
	}

	id, idErr := primitive.ObjectIDFromHex(project.ID.ToString())
	if idErr != nil {
		loggerObj.Error(idErr)
		return nil, idErr
	}

	if result, err := collection.UpdateOne(context.TODO(),
		bson.D{{utils.PROJECT_ID_FIELD, id}},
		projectDb); err != nil {
		loggerObj.Error(err)
		return nil, err
	} else {
		if result.MatchedCount != 1 {
			msg := fmt.Sprintf("Nothing to update. Project (%s) not found.", project.ID.ToString())
			loggerObj.Errorf(msg)
			return nil, errors.New(msg)
		}

		return project, nil
	}
}

func (dbManager *mongodbManagerImp) init() {
	// users collection
	if !dbManager.initiated {
		dbManager.cache = make(map[cacheKey]interface{})
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

		if result, err := collection.InsertOne(context.TODO(), userDB); err != nil {
			loggerObj.Error(err.Error())
			return nil, err
		} else {
			user.ID = &mdbId{id: result.InsertedID.(primitive.ObjectID)}
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

	dbId, err := dbManager.modelIDtoMongoID(id, loggerObj)

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
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if result.Err() != nil {
		loggerObj.Error(result.Err())
		return nil, result.Err()
	}

	userDb := &mdbUserModel{}
	if err := result.Decode(userDb); err != nil {
		loggerObj.Error(err)
		return nil, err
	}

	return userDb.toModel(), nil
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

	return dbManager.decodeBsonIntoUserListModel(cursor, loggerObj)
}

func (dbManager *mongodbManagerImp) decodeBsonIntoUserListModel(cursor *mongo.Cursor, loggerObj *logger.Logger) (model.UserList, error) {
	var users model.UserList

	for cursor.Next(context.TODO()) {
		userDb := &mdbUserModel{}
		if err := cursor.Decode(&userDb); err != nil {
			loggerObj.Error(err.Error())
			return nil, err
		}

		user := userDb.toModel()
		users = append(users, user)
	}

	return users, nil
}

func (dbManager *mongodbManagerImp) mongoIdToModelID(id primitive.ObjectID) model.ID {
	return &mdbId{id: id}
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
