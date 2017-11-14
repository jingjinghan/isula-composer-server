package storage

import (
	"errors"

	"github.com/astaxie/beego/logs"

	"github.com/isula/ihub/config"
	"github.com/isula/ihub/storage/driver"
)

var (
	sysDriver driver.StorageDriver
)

func init() {
	sysDriver = nil
}

// TODO: more logs
func loadDriver(cfg config.StorageConfig) (driver.StorageDriver, error) {
	for n, paras := range cfg {
		logs.Debug("Find storage driver for: %s", n)
		d, err := driver.FindDriver(n, paras)
		if err == nil {
			// Pickup the first qualified driver
			err = d.Init(paras)
			if err == nil {
				return d, nil
			}
		}
	}

	return nil, errors.New("Fail to get a suitable storage driver")
}

// InitStorage setups the storage bankends from the config file
func InitStorage(cfg config.StorageConfig) error {
	var err error
	sysDriver, err = loadDriver(cfg)
	// TODO: we should check the healthy status at the beginning

	return err
}

// Driver returns the storage driver
func Driver() driver.StorageDriver {
	cfg := config.GetConfig()
	if cfg.StorageLoad == "static" && sysDriver != nil {
		return sysDriver
	}

	var err error
	sysDriver, err = loadDriver(cfg.Storage)
	if err != nil {
		panic("Failed to load driver")
	}

	return sysDriver
}
