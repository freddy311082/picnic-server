package dbmanager

import (
	"github.com/freddy311082/picnic-server/model"
)

type DBManager interface {
	Open() error
	Close() error
	IsOpen() bool
	RegisterNewUser(user *model.User) (*model.User, error)
	GetUser(email string) (*model.User, error)
	AllUsers(startPosition, offset int) (model.UserList, error)
}

var dbManagerInstance DBManager

func DBManagerInstance() DBManager {

	ch := make(chan DBManager)

	go func() {
		if dbManagerInstance == nil {
			dbManagerInstance = createMongoDbManager()
		}

		ch <- dbManagerInstance
	}()

	return <-ch
}
