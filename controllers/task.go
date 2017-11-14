package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Task defines the task operations
type Task struct {
	beego.Controller
}

// Create adds a new task
func (t *Task) Create() {
	user := t.Ctx.Input.Param(":user")

	logs.Debug("Create a task for '%s'", user)
	CtxSuccessWrap(t.Ctx, http.StatusOK, "{}", nil)
}

// List lists the tasks
func (t *Task) List() {
	user := t.Ctx.Input.Param(":user")

	logs.Debug("List tasks of '%s'", user)
	CtxSuccessWrap(t.Ctx, http.StatusOK, "{}", nil)
}
