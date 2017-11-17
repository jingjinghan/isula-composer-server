package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/isula-composer-server/controllers"
)

const (
	taskPrefix = "/:user/task"
)

func init() {
	if err := RegisterRouter(taskPrefix, taskNameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", taskPrefix)
	} else {
		logs.Debug("Register router '%s' registered.", taskPrefix)
	}
}

// taskNameSpace defines the task router
func taskNameSpace() *beego.Namespace {
	ns := beego.NewNamespace(taskPrefix,
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),
		beego.NSRouter("/", &controllers.Task{}, "post:Create;get:List"),
		beego.NSRouter("/:id", &controllers.Task{}, "get:Get"),
	)

	return ns
}
