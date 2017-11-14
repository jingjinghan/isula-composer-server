package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/isula/isula-composer-server/controllers"
)

const (
	hookPrefix = "/:user/hook"
)

func init() {
	if err := RegisterRouter(hookPrefix, hookNameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", hookPrefix)
	} else {
		logs.Debug("Register router '%s' registered.", hookPrefix)
	}
}

// hookNameSpace defines the hook router
func hookNameSpace() *beego.Namespace {
	ns := beego.NewNamespace(hookPrefix,
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),
		beego.NSRouter("/", &controllers.Hook{}, "get:List;post:Create"),
		beego.NSRouter("/:id", &controllers.Hook{}, "get:Get;delete:Delete"),
	)

	return ns
}
