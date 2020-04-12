package dbmanager

import (
	"github.com/freddy311082/picnic-server/model"
)

type DBManager interface {
	Open() error
	Close() error
	IsOpen() bool
	RegisterNewUser(user *model.User) (*model.User, error)
	DeleteUser(email string) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id model.ID) (*model.User, error)
	AllUsers(startPosition, offset int) (model.UserList, error)
	AllUsersWhereIDIsIn(ids model.IDList) (model.UserList, error)
	AllProjects(startPosition, offset int) (model.ProjectList, error)
	AllProjectFromUser(user *model.User) (model.ProjectList, error)
	CreateProject(project *model.Project) (*model.Project, error)
	GetProject(projectId model.ID) (*model.Project, error)
	UpdateProject(project *model.Project) (*model.Project, error)
	DeleteProject(projectId model.ID) error
	DeleteProjects(ids model.IDList) error
	AllProjectWhereIDIsIn(ids model.IDList) (model.ProjectList, error)
	AddCustomer(customer *model.Customer) (*model.Customer, error)
	UpdateCustomer(customer *model.Customer) (*model.Customer, error)
	DeleteCustomer(customerId model.ID) error
	DeleteCustomers(ids model.IDList) error
	AllCustomers() (model.CustomerList, error)
	AllCustomersWhereIDIsIn(ids model.IDList) (model.CustomerList, error)
}

var dbManagerInstance DBManager

func Instance() DBManager {

	ch := make(chan DBManager)

	go func() {
		if dbManagerInstance == nil {
			dbManagerInstance = createMongoDbManager()
		}

		ch <- dbManagerInstance
	}()

	return <-ch
}
