package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"

	"github.com/isula/ihub/config"
	"github.com/isula/ihub/logger"
	_ "github.com/isula/ihub/storage/driver/filesystem"

	"github.com/isula/isula-composer-server/models"
	"github.com/isula/isula-composer-server/routers"
	"github.com/isula/isula-composer-server/storage"
)

func main() {
	cfg, err := config.InitConfigFromFile("conf/isula.yml")
	if err != nil {
		return
	}
	if err := logger.InitLogger(cfg.Log); err != nil {
		logs.Warning(err)
	}

	conn, _ := cfg.DB.GetConnection()
	if err := models.InitDB(conn, cfg.DB.Driver, "default"); err != nil {
		logs.Critical("Error in init db: ", err)
		return
	}

	if err := storage.InitStorage(cfg.Storage); err != nil {
		logs.Critical("Error in init storage: ", err)
		return
	}

	nss := router.GetNamespaces()
	for name, ns := range nss {
		logs.Debug("Namespace '%s' is enabled", name)
		beego.AddNamespace(ns)
	}

	beego.Run()
}
