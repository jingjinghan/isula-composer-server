package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	// TODO: makes error better name.

	// ErrEmptyDBUserOrPassword represents the 'user' or 'password' empty error
	ErrEmptyDBUserOrPassword = errors.New("database 'User' or 'Password' should not be empty")
	// ErrEmptyDBName represents the 'name' empty error
	ErrEmptyDBName = errors.New("database 'Name' should not be empty")
	// ErrEmptyDBServer represents the 'database server' empty error
	ErrEmptyDBServer = errors.New("database 'Server' should not be empty")

	// ErrNonStorageBackend represents the 'storage backend' empty error
	ErrNonStorageBackend = errors.New("at least one storage backend required")
	// ErrNonSessionBackend represents the 'session backend' empty error
	ErrNonSessionBackend = errors.New("at least one session backend required")
)

// Config defines the config items
type Config struct {
	// Port is the server listen port
	Port int64 `yaml:"port"`
	// Log defines the log config group
	// The log validation will be checked in the log init part.
	Log LogConfig `yaml:"log,omitempty"`
	// StorageLoad
	// 'dynamic' or 'static'
	// static: load at the first time
	// dynamic: load every time, most time because of multiple tenant using their own token/ak-sk
	// TODO: should have 'default' value
	StorageLoad string `yaml:"storageload,omitempty"`
	// Storage defines the storage config group
	Storage StorageConfig `yaml:"storage"`
	// DB defines the db config group
	DB DBConfig `yaml:"db"`
	// Session defines the session config group
	Session SessionConfig `yaml:"session"`
}

// Valid checks if a config is logical
func (cfg *Config) Valid() error {
	if err := cfg.Storage.Valid(); err != nil {
		return err
	}
	if err := cfg.DB.Valid(); err != nil {
		return err
	}

	return cfg.Session.Valid()
}

// LogConfig stores the log config item group
type LogConfig map[string](map[string]interface{})

// StorageConfig stores the storage config item group
type StorageConfig map[string](map[string]interface{})

// Valid checks the storage config validation
func (cfg *StorageConfig) Valid() error {
	if len(*cfg) == 0 {
		return ErrNonStorageBackend
	}

	return nil
}

// SessionConfig stores the session config item group
type SessionConfig map[string](map[string]interface{})

// Valid checks the session config validation
func (cfg *SessionConfig) Valid() error {
	if len(*cfg) == 0 {
		return ErrNonSessionBackend
	}

	return nil
}

// DBConfig defines the db configs
type DBConfig struct {
	Driver   string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Name     string `yaml:"name"`
}

// GetConnection returns the sql recognizable connection string
func (cfg *DBConfig) GetConnection() (string, error) {
	if cfg.User == "" || cfg.Password == "" {
		return "", ErrEmptyDBUserOrPassword
	}

	if cfg.Name == "" {
		return "", ErrEmptyDBName
	}
	if cfg.Server == "" {
		return "", ErrEmptyDBServer
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cfg.User, cfg.Password, cfg.Server, cfg.Name), nil
}

// Valid checks the db config validation
func (cfg *DBConfig) Valid() error {
	_, err := cfg.GetConnection()
	if err != nil {
		return err
	}

	return nil
}

var (
	sysConfig Config
)

// InitConfigFromFile loads the config from a file
func InitConfigFromFile(path string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	// TODO: add lock?
	sysConfig = config

	return config, nil
}

// GetConfig returns the current system config setting
func GetConfig() Config {
	return sysConfig
}
