package service

import (
	"errors"
	"github.com/freddy311082/picnic-server/dbmanager"
	"github.com/freddy311082/picnic-server/model"
	"github.com/freddy311082/picnic-server/utils"
)

type privateId struct {
	id string
}

func (objId *privateId) ToString() string {
	return objId.id
}

type Service interface {
	Init() error
	AllUsers(startPosition, offset int) (model.UserList, error)
	RegisterUser(user *model.User) (*model.User, error)
	GetUser(user *model.User) (*model.User, error)
	DeleteUser(user *model.User) error
	CreateProject(project *model.Project) (*model.Project, error)
	AllProjects(startPosition, offset int) (model.ProjectList, error)
	AllProjectsByUser(user *model.User) (model.ProjectList, error)
	AddProject(project *model.Project) (*model.Project, error)
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
	NewIDFromString(strId string) model.ID
}

type serviceImp struct {
	dbManager dbmanager.DBManager
}

func (service *serviceImp) AllProjectWhereIDIsIn(ids model.IDList) (model.ProjectList, error) {
	return dbmanager.Instance().AllProjectWhereIDIsIn(ids)
}

func (service *serviceImp) AddCustomer(customer *model.Customer) (*model.Customer, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if customer != nil || customer.Name == "" {
		const msg = "unable to create a customer. Customer cannot be null or name cannot be empty"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	return dbmanager.Instance().AddCustomer(customer)
}

func (service *serviceImp) UpdateCustomer(customer *model.Customer) (*model.Customer, error) {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()

	if customer != nil || customer.Name == "" {
		const msg = "unable to create a customer. Customer cannot be null or name cannot be empty"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	return dbmanager.Instance().UpdateCustomer(customer)
}

func (service *serviceImp) DeleteCustomer(customerId model.ID) error {
	return dbmanager.Instance().DeleteCustomer(customerId)
}

func (service *serviceImp) DeleteCustomers(ids model.IDList) error {
	return dbmanager.Instance().DeleteCustomers(ids)
}

func (service *serviceImp) AllCustomers() (model.CustomerList, error) {
	return dbmanager.Instance().AllCustomers()
}

func (service *serviceImp) AllCustomersWhereIDIsIn(ids model.IDList) (model.CustomerList, error) {
	return dbmanager.Instance().AllCustomersWhereIDIsIn(ids)
}

func (service *serviceImp) GetUser(user *model.User) (*model.User, error) {
	if user == nil {
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		const msg = "user object cannot be null"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	if user.ID != nil {
		return dbmanager.Instance().GetUserByID(user.ID)
	} else {
		return dbmanager.Instance().GetUserByEmail(user.Email)
	}
}

func (service *serviceImp) DeleteUser(user *model.User) error {
	if user == nil {
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		loggerObj.Error()
	}

	return dbmanager.Instance().DeleteUser(user.Email)
}

func (service *serviceImp) CreateProject(project *model.Project) (*model.Project, error) {
	if project == nil {
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		const msg = "project object cannot be null"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	return dbmanager.Instance().CreateProject(project)
}

func (service *serviceImp) AllProjects(startPosition, offset int) (model.ProjectList, error) {
	return dbmanager.Instance().AllProjects(startPosition, offset)
}

func (service *serviceImp) AllProjectsByUser(user *model.User) (model.ProjectList, error) {
	if user == nil {
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		const msg = "user object cannot be null"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	return dbmanager.Instance().AllProjectFromUser(user)
}

func (service *serviceImp) AddProject(project *model.Project) (*model.Project, error) {
	if project == nil {
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		const msg = "project object cannot be null"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	return dbmanager.Instance().CreateProject(project)
}

func (service *serviceImp) UpdateProject(project *model.Project) (*model.Project, error) {
	if project == nil {
		loggerObj := utils.LoggerObj()
		defer loggerObj.Close()
		const msg = "project object cannot be null"
		loggerObj.Error(msg)
		return nil, errors.New(msg)
	}

	return dbmanager.Instance().UpdateProject(project)
}

func (service *serviceImp) DeleteProject(projectId model.ID) error {
	return dbmanager.Instance().DeleteProject(projectId)
}

func (service *serviceImp) DeleteProjects(ids model.IDList) error {
	return dbmanager.Instance().DeleteProjects(ids)
}

func (service *serviceImp) Init() error {
	loggerObj := utils.LoggerObj()
	defer loggerObj.Close()
	if err := serviceInstance.dbManager.Open(); err != nil {
		loggerObj.Error(err.Error())
		panic(err.Error())
		return err
	}

	return nil
}

func (service *serviceImp) NewIDFromString(strId string) model.ID {
	return &privateId{id: strId}
}

func (service *serviceImp) RegisterUser(user *model.User) (*model.User, error) {
	return dbmanager.Instance().RegisterNewUser(user)
}

func (service *serviceImp) AllUsers(startPosition, offset int) (model.UserList, error) {
	return dbmanager.Instance().AllUsers(startPosition, offset)
}

var serviceInstance *serviceImp

func Instance() Service {
	ch := make(chan Service)

	go func() {
		if serviceInstance == nil {
			serviceInstance = &serviceImp{
				dbManager: dbmanager.Instance(),
			}
		}

		ch <- serviceInstance
	}()

	return <-ch
}
