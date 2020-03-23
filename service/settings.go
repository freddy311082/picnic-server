package service

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/logger"
	"io/ioutil"
)

const CONFIG_FILE_PATH = "./config/settings.json"

// ******************************* Interfaces ***********************************

type DBSettings interface {
	Host() string
	Port() uint
	DbName() string
	User() string
	Password() string
}

type Settings interface {
	DBSettingsValues() DBSettings
}

// ******************************* Struct ***********************************

// ******************************* dbSettingsImp ***********************************

type dbSettingsImp struct {
	_host     string
	_port     uint
	_dbName   string
	_user     string
	_password string
}

func (dbSettings *dbSettingsImp) Host() string {
	return dbSettings._host
}

func (dbSettings *dbSettingsImp) Port() uint {
	return dbSettings._port
}

func (dbSettings *dbSettingsImp) DbName() string {
	return dbSettings._dbName
}

func (dbSettings *dbSettingsImp) User() string {
	return dbSettings._user
}

func (dbSettings *dbSettingsImp) Password() string {
	return dbSettings._password
}

func (dbSettings *dbSettingsImp) loadData(data map[string]interface{}) error {
	dbDataSettings := data["db"].(map[string]interface{})

	dbSettings._host = dbDataSettings["host"].(string)
	dbSettings._port = dbDataSettings["port"].(uint)
	dbSettings._dbName = dbDataSettings["dbname"].(string)
	dbSettings._user = dbDataSettings["user"].(string)
	dbSettings._password = dbDataSettings["password"].(string)

	var err error

	if err = dbSettings.validate(); err == nil {
		dbSettings.encryptPassword()
	}

	return err
}

func (dbSettings *dbSettingsImp) validate() error {
	if (len(dbSettings._host) == 0 &&
		dbSettings._port <= 1024 &&
		len(dbSettings._dbName) == 0 &&
		len(dbSettings._user) == 0) &&
		len(dbSettings._password) == 0 {
		const msg = "invalid database settings values. Please check and rerun the server again"
		logger.Error(msg)
		logger.Info(fmt.Sprint(`host: %s
port: %d
dbname: %s
user: %s
password: %s`, dbSettings._host, dbSettings._port, dbSettings._dbName, dbSettings._user, dbSettings._password))
		return errors.New(msg)
	}

	return nil
}

func (dbSettings *dbSettingsImp) encryptPassword() {
	newPassword := sha256.Sum256(([]byte)("Picnic:" + dbSettings._password))
	dbSettings._password = fmt.Sprint(newPassword)
}

// ******************************* dbSettingsImp ***********************************

type settingsImp struct {
	dbSettings *dbSettingsImp
}

func (settings *settingsImp) DBSettingsValues() DBSettings {
	return settings.dbSettings
}

func (settings *settingsImp) load() error {
	var err error
	var content []byte
	if content, err = ioutil.ReadFile(CONFIG_FILE_PATH); err != nil {
		return err
	}

	return settings.loadContent(content)
}

func (settings *settingsImp) loadContent(content []byte) error {
	if content == nil || len(content) == 0 {
		msg := "invalid setting.json file content"
		logger.Error(msg)
		return errors.New(msg)
	} else {

		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			const msg = "error parsing content of settings.json"
			logger.Error(msg)
			return errors.New(msg)
		}
		if err := settings.loadDbSettings(data); err != nil {
			return err
		}
	}

	return nil
}

func (settings *settingsImp) loadDbSettings(data map[string]interface{}) error {
	settings.dbSettings = &dbSettingsImp{}
	return settings.dbSettings.loadData(data)
}

// ******************************* Public Functions ***********************************

var settingsSingleton *settingsImp
