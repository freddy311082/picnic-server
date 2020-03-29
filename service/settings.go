package service

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/freddy311082/picnic-server/utils"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

const CONFIG_FILE_PATH = "./config/settings.json"

// ******************************* Interfaces ***********************************

type DBSettings interface {
	Host() string
	Port() int
	DbName() string
	User() string
	Password() string
	DriverType() utils.DBTypeEnum

	ChangeDatabase(dbName string)
	ConnectionString() string
	ToString() string
}

type APISettings interface {
	GraphiQL() bool
	HttpPort() int
	ToString() string
}

type Settings interface {
	DBSettingsValues() DBSettings
	APISettings() APISettings
	ToString() string
	filename() string
}

// ******************************* Struct ***********************************

// ******************************* apiSettingsImp ***********************************

type apiSettingsImp struct {
	allowGraphiQL bool
	httpPort      int
}

func (apiSettings *apiSettingsImp) ToString() string {
	return fmt.Sprintf(`
========== API Settings =========
Allowed GraphiQL: %s
HTTP Port: %d
=================================

`, fmt.Sprint(apiSettings.allowGraphiQL), apiSettings.httpPort)
}

func (apiSettings *apiSettingsImp) HttpPort() int {
	return apiSettings.httpPort
}

func (apiSettings *apiSettingsImp) GraphiQL() bool {
	return apiSettings.allowGraphiQL
}

func (apiSettings *apiSettingsImp) loadData(data map[string]interface{}) error {
	if apiMap, ok := data[utils.WEBSERVER_JSON_KEY].(map[string]interface{}); !ok {
		const msg = "invalid settings.json. Error or missing API config"
		utils.PicnicLog_ERROR(msg)
		return errors.New(msg)
	} else {
		apiSettings.allowGraphiQL = apiMap[utils.GRAPHIQL_JSON_KEY].(bool)
		apiSettings.httpPort = int(apiMap[utils.HTTP_PORT_JSON_KEY].(float64))
	}

	return nil
}

// ******************************* dbSettingsImp ***********************************

type dbSettingsImp struct {
	_host     string
	_port     int
	_dbName   string
	_user     string
	_password string
}

func (dbSettings *dbSettingsImp) ToString() string {
	return fmt.Sprintf(`
========== Database Settings =========
Host: %s
DB Port: %d
DB Name: %s
DB User: %s
DB Password: %s
Connection String: %s
=================================
`, dbSettings._host, dbSettings._port, dbSettings._dbName, dbSettings._user, dbSettings._password,
		dbSettings.ConnectionString())
}

func (dbSettings *dbSettingsImp) Host() string {
	return dbSettings._host
}

func (dbSettings *dbSettingsImp) Port() int {
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

func (dbSettings *dbSettingsImp) DriverType() utils.DBTypeEnum {
	return utils.DBType_MONGODB
}

func (dbSettings *dbSettingsImp) ChangeDatabase(dbName string) {
	dbSettings._dbName = dbName
}

func (dbSettings *dbSettingsImp) ConnectionString() string {
	hostWithDb := fmt.Sprintf(dbSettings._host, dbSettings._dbName)
	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s", dbSettings._user, dbSettings._password, hostWithDb)
	return connectionString
}

func (dbSettings *dbSettingsImp) loadData(data map[string]interface{}) error {
	if dbMap, okDb := data["db"].(map[string]interface{}); !okDb {
		return errors.New("missing \"db\" key in settings.json")
	} else if mongodbMap, okMongodb := dbMap["mongodb"].(map[string]interface{}); !okMongodb {
		return errors.New("missing \"mongodb\" key in settings.json")
	} else {
		dbSettings._host = mongodbMap["host"].(string)
		dbSettings._port = int(mongodbMap["port"].(float64))
		dbSettings._dbName = mongodbMap["dbname"].(string)
		dbSettings._user = mongodbMap["user"].(string)
		dbSettings._password = mongodbMap["password"].(string)
	}

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
		utils.PicnicLog_ERROR(msg)
		utils.PicnicLog_INFO(fmt.Sprintf(`host: %s
port: %d
dbname: %s
user: %s
password: %s`, dbSettings._host, dbSettings._port, dbSettings._dbName, dbSettings._user, dbSettings._password))
		return errors.New(msg)
	}

	return nil
}

func (dbSettings *dbSettingsImp) encryptPassword() {
	hash := sha256.New()
	hash.Write(([]byte)(dbSettings._password))

	dbSettings._password = base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// ******************************* dbSettingsImp ***********************************

type settingsImp struct {
	dbSettings  *dbSettingsImp
	apiSettings *apiSettingsImp
}

func (settings *settingsImp) APISettings() APISettings {
	return settings.apiSettings
}

func (settings *settingsImp) DBSettingsValues() DBSettings {
	return settings.dbSettings
}

func (settings *settingsImp) ToString() string {
	return settings.dbSettings.ToString() + settings.apiSettings.ToString()
}

func (settings *settingsImp) filename() string {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		panic("settings.json file not found")
		return ""
	}

	baseDir := path.Dir(filename)
	result := baseDir + string(os.PathSeparator) + ".." + string(os.PathSeparator) + path.Join("config", "settings.json")

	return result
}

func (settings *settingsImp) fileContent() ([]byte, error) {

	filename := settings.filename()
	var content, err = ioutil.ReadFile(filename)
	if err != nil {
		msg := fmt.Sprintf("Error reading file %s", filename)
		utils.PicnicLog_ERROR(msg)
		return nil, err
	}

	return content, nil
}

func (settings *settingsImp) load() error {
	if content, err := settings.fileContent(); err != nil {
		return err
	} else {
		return settings.loadContent(content)
	}
}

func (settings *settingsImp) loadContent(content []byte) error {
	if content == nil || len(content) == 0 {
		msg := "invalid setting.json file content"
		utils.PicnicLog_ERROR(msg)
		return errors.New(msg)
	} else {

		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			const msg = "error parsing content of settings.json"
			utils.PicnicLog_ERROR(msg)
			return errors.New(msg)
		}
		if err := settings.loadDbSettings(data); err != nil {
			return err
		} else if err = settings.loadApiSettings(data); err != nil {
			return err
		}
	}

	return nil
}

func (settings *settingsImp) loadDbSettings(data map[string]interface{}) error {
	settings.dbSettings = &dbSettingsImp{}
	return settings.dbSettings.loadData(data)
}

func (settings *settingsImp) loadApiSettings(data map[string]interface{}) error {
	settings.apiSettings = &apiSettingsImp{}
	return settings.apiSettings.loadData(data)
}

// ******************************* Public Functions ***********************************

var settingsSingleton *settingsImp

func SettingsObj() Settings {

	type result struct {
		instance *settingsImp
		err      error
	}

	ch := make(chan result)

	go func() {
		if settingsSingleton == nil {
			settingsSingleton = &settingsImp{}
			err := settingsSingleton.load()
			ch <- result{
				instance: settingsSingleton,
				err:      err,
			}
		}

		ch <- result{
			instance: settingsSingleton,
			err:      nil,
		}
	}()

	value := <-ch
	return value.instance
}
