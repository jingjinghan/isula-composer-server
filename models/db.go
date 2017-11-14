package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// InitDB regists the orm db driver and regists to the database
func InitDB(conn string, driver string, name string) error {
	if driver != "mysql" {
		return errors.New("only support mysql yet")
	}

	orm.RegisterDriver(driver, orm.DRMySQL)
	orm.DefaultTimeLoc = time.UTC

	for i := 0; i < 3; i++ {
		// TODO: do we need this?
		//		<-time.After(time.Second * 5)
		err := orm.RegisterDataBase(name, driver, conn)
		if err == nil {
			orm.SetMaxIdleConns(name, 30)
			orm.SetMaxOpenConns(name, 30)
			orm.RunSyncdb(name, false, true)
			return nil
		}
		logs.Debug("Try to register database for %d times...", i)
	}

	return errors.New("Cannot connect to database")
}
