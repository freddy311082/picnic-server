package dbmanager

import (
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

func (dbUser *mdbUserModel) initFromModel(user *model.User) {
	if user.ID != nil {
		if objId, err := primitive.ObjectIDFromHex(user.ID.ToString()); err == nil {
			dbUser.ID = objId
		}
	}
	dbUser.Name = user.Name
	dbUser.LastName = user.LastName
	dbUser.Email = user.Email
}

func (dbUser *mdbUserModel) toModel() *model.User {
	return &model.User{
		ID: &mdbId{
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
	CustomerID  primitive.ObjectID `bson:"customer_id"`
}

func (dbProject *mdbProjectModel) initFromModel(project *model.Project) error {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if project.ID != nil {
		if objId, err := primitive.ObjectIDFromHex(project.ID.ToString()); err != nil {
			dbProject.ID = objId
		}
	}

	if id, err := primitive.ObjectIDFromHex(project.Owner.ID.ToString()); err != nil {
		loggerObj.Error(err)
		return err
	} else {
		dbProject.OwnerID = id
	}

	if id, err := primitive.ObjectIDFromHex(project.Customer.ID.ToString()); err != nil {
		loggerObj.Error(err)
		return err
	} else {
		dbProject.CustomerID = id
	}

	dbProject.Name = project.Name
	dbProject.Description = project.Description
	dbProject.CreatedAt = primitive.NewDateTimeFromTime(project.CreatedAt)

	return nil
}

func (dbProject *mdbProjectModel) toModel() *model.Project {
	var ownerId *mdbId
	var customerId *mdbId

	if dbProject.CustomerID.Hex() != "" {
		ownerId = &mdbId{id: dbProject.CustomerID}
	}

	if dbProject.OwnerID.Hex() != "" {
		customerId = &mdbId{id: dbProject.OwnerID}
	}

	return &model.Project{
		ID:          &mdbId{id: dbProject.ID},
		Name:        dbProject.Name,
		Description: dbProject.Description,
		CreatedAt:   dbProject.CreatedAt.Time(),
		Owner:       &model.User{ID: customerId},
		Customer:    &model.Customer{ID: ownerId},
	}
}

type mdbProjectListModel []mdbProjectModel

func (dbProjectList mdbProjectListModel) toModel() model.ProjectList {
	result := model.ProjectList{}

	for _, dbProject := range dbProjectList {
		result = append(result, dbProject.toModel())
	}

	return result
}

type mdbCustomerModel struct {
	ID       primitive.ObjectID   `bson:"_id"`
	Name     string               `bson:"name"`
	Cuit     string               `bson:"cuit"`
	Projects []primitive.ObjectID `bson:"projects"`
}

func (dbCustomer *mdbCustomerModel) toModel() (*model.Customer, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	customer := &model.Customer{
		ID:   &mdbId{id: dbCustomer.ID},
		Name: dbCustomer.Name,
		Cuit: dbCustomer.Cuit,
	}

	dbManager := Instance().(*mongodbManagerImp)
	ids := dbManager.mongoIDsToModelIDs(dbCustomer.Projects)
	projects, err := dbManager.AllProjectWhereIDIsIn(ids)
	if err == nil {
		customer.Projects = projects
	}

	return customer, err
}

func (dbCustomer *mdbCustomerModel) initFromModel(customer *model.Customer) {
	if customer.ID != nil {
		if objId, err := primitive.ObjectIDFromHex(customer.ID.ToString()); err == nil {
			dbCustomer.ID = objId
		}
	}

	dbCustomer.Name = customer.Name
	dbCustomer.Cuit = customer.Cuit

	dbManager := Instance().(*mongodbManagerImp)
	dbCustomer.Projects, _ = dbManager.modelIDsToMongoIDs(customer.Projects.IDs(), utils.LoggerObj())
}
