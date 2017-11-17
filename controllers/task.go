package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/isula/isula-composer-server/models"
)

// Task defines the task operations
type Task struct {
	beego.Controller
}

// Create adds a new task
func (t *Task) Create() {
	user := t.Ctx.Input.Param(":user")
	output := t.Ctx.Input.Query("output")

	logs.Debug("Create task for '%s', output is '%s'", user, output)

	var config, scripts string

	file, _, err := t.Ctx.Request.FormFile("file")
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusBadRequest, nil, fmt.Sprintf("Cannot find the upload file '%s'.", user))
		return
	}
	data, _ := ioutil.ReadAll(file)
	file.Close()

	if output != "" {
		config = string(data)
	} else {
		scripts = string(data)
	}

	id, err := models.AddTaskFull(user, output, config, scripts)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to create a task for '%s'.", user))
		return
	}

	CtxSuccessWrap(t.Ctx, http.StatusOK, id, nil)
}

// List lists the tasks
func (t *Task) List() {
	user := t.Ctx.Input.Param(":user")

	logs.Debug("List tasks of '%s'", user)
	CtxSuccessWrap(t.Ctx, http.StatusOK, "{}", nil)
}

// Get returns the task detail
func (t *Task) Get() {
	user := t.Ctx.Input.Param(":user")
	idStr := t.Ctx.Input.Param(":id")

	logs.Debug("Get task %s from '%s'", idStr, user)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusBadRequest, err, fmt.Sprintf("Invalid id detected '%s': %v.", idStr, err))
		return
	}

	task, err := models.QueryTaskByID(int64(id))
	if err != nil {
		CtxErrorWrap(t.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get the task '%s' from '%s'.", id, user))
		return
	} else if task == nil {
		CtxErrorWrap(t.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to find the task '%s' from '%s'.", id, user))
		return
	}

	CtxSuccessWrap(t.Ctx, http.StatusOK, task, nil)
}
